package user

import (
	"context"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/controller/contextextensions"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
)

func (u *User) Delete(ctx context.Context, args *struct{}, vars *messages.RouteVars) (*struct{}, *wrapper.HandlerError) {
	claims, err := contextextensions.ClaimsFromCtx(ctx)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	userID, err := validateUpdateAccess(claims, nil, vars)
	if err != nil || userID == nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusNotFound}
	}

	err = u.dal.User.DeleteUser(ctx, *userID, claims.AccessLevel)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusNotFound}
	}

	return &struct{}{}, nil
}
