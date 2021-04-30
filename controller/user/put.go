package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/controller/contextextensions"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
)

func (u *User) Put(ctx context.Context, args *messages.UpdateUserRequest, vars *struct{}) (*messages.BaseResponse, *wrapper.HandlerError) {
	claims, err := contextextensions.ClaimsFromCtx(ctx)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	user, err := u.dal.User.GetUserByID(ctx, claims.UserID, nil)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusNotFound}
	}

	err = updateUser(user, args, claims.AccessLevel)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusUnauthorized}
	}

	err = u.dal.User.UpsertUser(ctx, user)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	return &messages.BaseResponse{Code: http.StatusOK}, nil
}

func updateUser(user *models.User, args *messages.UpdateUserRequest, accessLevel models.AccessLevel) error {
	if args.CalorieLimit != nil {
		user.CalorieLimit = *args.CalorieLimit
	}
	if args.Username != nil {
		user.UserName = *args.Username
	}

	if args.AccessLevel != nil {
		if accessLevel != models.Admin {
			return fmt.Errorf("only admin can change access level")
		}
	}

	if args.Password != nil {
		newPassword, err := contextextensions.HashPassword(*args.Password)
		if err != nil {
			return err
		}
		user.Password = newPassword
	}
	return nil
}
