package controller

import (
	"context"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/dal"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/spf13/viper"
)

type User struct {
	dal    *dal.DAL
	config *viper.Viper
}

func NewUser(
	dal *dal.DAL,
	config *viper.Viper,
) *User {
	return &User{
		dal:    dal,
		config: config,
	}
}

func (u *User) Put(ctx context.Context, args *messages.UpdateUserRequest, vars *struct{}) (*messages.BaseResponse, error) {
	claims, err := ClaimsFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	user, err := u.dal.User.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	updateUser(user, args)

	err = u.dal.User.UpsertUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &messages.BaseResponse{Code: http.StatusOK}, nil
}

func updateUser(user *models.User, args *messages.UpdateUserRequest) {
	if args.CalorieLimit != nil {
		user.CalorieLimit = *args.CalorieLimit
	}
	if args.Username != nil {
		user.UserName = *args.Username
	}
}
