package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/controller/contextextensions"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
)

func (u *User) Put(ctx context.Context, args *messages.UpdateUserRequest, vars *messages.RouteVars) (*struct{}, *wrapper.HandlerError) {
	claims, err := contextextensions.ClaimsFromCtx(ctx)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	userID, err := validateUpdateAccess(claims, args, vars)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusUnauthorized}
	}

	user, err := u.dal.User.GetUsers(ctx, userID, nil)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusNotFound}
	}

	err = updateUser(user[0], args, claims.AccessLevel)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusUnauthorized}
	}

	err = u.dal.User.UpsertUser(ctx, user[0], claims.AccessLevel)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusNotFound}
	}

	return &struct{}{}, nil
}

func validateUpdateAccess(claims *models.Claims, args *messages.UpdateUserRequest, vars *messages.RouteVars) (*uuid.UUID, error) {
	noUserIDSpecified := (vars == nil || vars.UserID == nil)
	if noUserIDSpecified {
		return nil, fmt.Errorf("no user specified in put user")
	}

	if *vars.UserID == "me" {
		return &claims.UserID, nil
	}

	userID, err := uuid.Parse(*vars.UserID)
	if err != nil {
		return nil, err
	}

	if userID == claims.UserID {
		return &claims.UserID, nil
	}

	hasAccess := claims.AccessLevel == models.Admin || claims.AccessLevel == models.Manager
	if hasAccess {
		return &userID, nil
	}
	return nil, fmt.Errorf("wrong acccess level")
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
		user.AccessLevel = *args.AccessLevel
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
