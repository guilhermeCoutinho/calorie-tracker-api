package http

import (
	"net/http"
	"strings"

	"github.com/guilhermeCoutinho/api-studies/usecase"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	usecase *usecase.Usecase
	logger  logrus.FieldLogger
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// format: Authorization: Bearer
		tokens, ok := r.Header["Authorization"]
		if ok && len(tokens) >= 1 {
			token = tokens[0]
			token = strings.TrimPrefix(token, "Bearer ")
		}

		if token == "" {
			m.logger.Error("Token is eempty")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		claims, err := m.usecase.ClaimsFromToken(token)
		if err != nil {
			m.logger.WithError(err).Error("Failed to retrieve uuid from token")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		logger := m.logger.WithFields(logrus.Fields{
			"userID":      claims.UserID,
			"accessLevel": claims.AccessLevel,
		})

		logger.Info("User authorized")
		r = r.WithContext(m.usecase.UUIDToCtx(r.Context(), claims.UserID))
		next.ServeHTTP(w, r)
	})
}
