package controller

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Auth struct {
	logger logrus.FieldLogger
}

func NewAuth(
	logger logrus.FieldLogger,
) *Auth {
	return &Auth{
		logger: logger,
	}
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *Auth) Authenticate(w http.ResponseWriter, r *http.Request) {
	args := &AuthRequest{}
	err := json.NewDecoder(r.Body).Decode(args)
	if err != nil {
		a.logger.WithError(err).Error("Failed to parse signup payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//usecase.Auth()
}
