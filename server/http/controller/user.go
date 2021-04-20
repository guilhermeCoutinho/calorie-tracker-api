package controller

import (
	"encoding/json"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/usecase"
	"github.com/sirupsen/logrus"
)

type User struct {
	logger logrus.FieldLogger
}

func NewUser(
	logger logrus.FieldLogger,
) *User {
	return &User{
		logger: logger,
	}
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	args := &AuthRequest{}
	err := json.NewDecoder(r.Body).Decode(args)
	if err != nil {
		u.logger.WithError(err).Error("Failed to parse signup payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = usecase.CreateUser(args.Username, args.Password)
	if err != nil {
		u.logger.WithError(err).Error("Failed to create user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (u *User) GetProfile(w http.ResponseWriter, r *http.Request) {
	// check token

	// return profile
	w.Write([]byte("Not implemented\n"))
}
