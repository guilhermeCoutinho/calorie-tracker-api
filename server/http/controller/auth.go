package controller

import (
	"encoding/json"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/usecase"
	"github.com/sirupsen/logrus"
)

type Auth struct {
	usecase *usecase.Usecase
	logger  logrus.FieldLogger
}

func NewAuth(
	usecase *usecase.Usecase,
	logger logrus.FieldLogger,
) *Auth {
	return &Auth{
		usecase: usecase,
		logger:  logger,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"token"`
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	args := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(args)
	if err != nil {
		a.logger.WithError(err).Error("Failed to parse signup payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := a.usecase.Login(r.Context(), args.Username, args.Password)
	if err != nil {
		a.logger.WithError(err).Error("Failed to create user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		AccessToken: token,
	}
	writeResponse(response, w)
	a.logger.Info("User login successfully")
}
