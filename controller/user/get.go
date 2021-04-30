package user

import (
	"context"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/controller/contextextensions"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
)

func (u *User) Get(ctx context.Context, args *struct{}, vars *messages.RouteVars) (*messages.GetUsersResponse, *wrapper.HandlerError) {
	claims, err := contextextensions.ClaimsFromCtx(ctx)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	canAccessOtherUsers := func(al models.AccessLevel) bool { return al == models.Admin || al == models.Manager }
	userID, wrapperErr := contextextensions.ValidatetGetAccessLevel(claims, vars, canAccessOtherUsers)
	if wrapperErr != nil {
		return nil, wrapperErr
	}

	user, err := u.dal.User.GetUserByID(ctx, *userID, contextextensions.GetQueryOptions(ctx))
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusNotFound}
	}

	return &messages.GetUsersResponse{
		Users:        user,
		BaseResponse: messages.BaseResponse{Code: http.StatusOK},
	}, nil
}
