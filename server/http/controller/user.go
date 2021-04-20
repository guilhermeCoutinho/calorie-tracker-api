package controller

import (
	"encoding/json"
	"net/http"

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

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	args := &CreateRequest{}
	err := json.NewDecoder(r.Body).Decode(args)
	if err != nil {
		u.logger.WithError(err).Error("Failed to parse signup payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = u.usecase.CreateUser(args.Username, args.Password)
	if err != nil {
		u.logger.WithError(err).Error("Failed to create user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	u.logger.Info("Player created successfully")
}

func (u *User) GetProfile(w http.ResponseWriter, r *http.Request) {
	// check token
	// return profile
	w.Write([]byte("Not implemented\n"))
}
