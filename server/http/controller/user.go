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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	u.logger.Info("User created successfully")
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"token"`
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	args := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(args)
	if err != nil {
		u.logger.WithError(err).Error("Failed to parse signup payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := u.usecase.UserLogin(args.Username, args.Password)
	if err != nil {
		u.logger.WithError(err).Error("Failed to create user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		AccessToken: user.Token,
	}
	writeResponse(response, w)
	u.logger.Info("User login successfully")
}

func writeResponse(data interface{}, w http.ResponseWriter) {
	bytes, _ := json.Marshal(data)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
