package meal

import (
	"context"
	"net/http"

	"github.com/guilhermeCoutinho/api-studies/controller/contextextensions"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
)

func (m *Meal) Delete(ctx context.Context, args *messages.DeleteMealRequest, vars *messages.RouteVars) (*struct{}, *wrapper.HandlerError) {
	claims, err := contextextensions.ClaimsFromCtx(ctx)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	accessFilter := func(al models.AccessLevel) bool { return al == models.Admin }
	userID, wrapperErr := contextextensions.ValidatetGetAccessLevel(claims, vars, accessFilter)
	if wrapperErr != nil {
		return nil, wrapperErr
	}

	mealID := m.getMealIDFromURL(&args.ID, vars)
	err = m.dal.Meal.DeleteMeal(ctx, *mealID, userID)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusNotFound}
	}

	return nil, nil
}
