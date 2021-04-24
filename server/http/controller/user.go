package controller

import (
	"encoding/json"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/usecase"
	"github.com/sirupsen/logrus"
)

type User struct {
	logger  logrus.FieldLogger
	usecase *usecase.Usecase
}

func NewUser(
	logger logrus.FieldLogger,
	usecase *usecase.Usecase,
) *User {
	return &User{
		logger:  logger,
		usecase: usecase,
	}
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	args := &messages.CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(args)
	if err != nil {
		u.logger.WithError(err).Error("Failed to parse signup payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = u.usecase.CreateUser(r.Context(), args.Username, args.Password)
	if err != nil {
		u.logger.WithError(err).Error("Failed to create user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	u.logger.Info("User created successfully")
}

func writeResponse(data interface{}, w http.ResponseWriter) {
	bytes, _ := json.Marshal(data)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
