package meal

import (
	"context"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/controller/contextextensions"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
)

func (m *Meal) Get(ctx context.Context, args *struct{}, vars *messages.RouteVars) (*messages.GetMealsResponse, *wrapper.HandlerError) {
	claims, err := contextextensions.ClaimsFromCtx(ctx)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	accessFilter := func(al models.AccessLevel) bool { return al == models.Admin }
	userID, wrapperErr := contextextensions.ValidatetGetAccessLevel(claims, vars, accessFilter)
	if wrapperErr != nil {
		return nil, wrapperErr
	}

	meals, err := m.dal.Meal.GetMeals(ctx, userID, contextextensions.GetQueryOptions(ctx))
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	return &messages.GetMealsResponse{
		Meals: meals,
	}, nil
}
