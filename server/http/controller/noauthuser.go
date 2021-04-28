package controller

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/dal"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type UserNoAuth struct {
	dal    *dal.DAL
	config *viper.Viper
}

func NewUserNoAuth(
	dal *dal.DAL,
	config *viper.Viper,
) *UserNoAuth {
	return &UserNoAuth{
		dal:    dal,
		config: config,
	}
}

func (u *UserNoAuth) Post(ctx context.Context, args *messages.CreateUserRequest, vars *struct{}) (*messages.BaseResponse, error) {
	hashedPassword, err := hashPassword(args.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		UserName:     args.Username,
		Password:     hashedPassword,
		CalorieLimit: args.CalorieLimit,
	}

	err = u.dal.User.UpsertUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &messages.BaseResponse{Code: http.StatusOK}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
