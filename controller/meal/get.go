package meal

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

func (m *Meal) Get(ctx context.Context, args *messages.GetMealsResponse, vars *messages.RouteVars) (*messages.GetMealsResponse, *wrapper.HandlerError) {
	claims, err := contextextensions.ClaimsFromCtx(ctx)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	userID, wrapperErr := m.validatetGetAccessLevel(claims, vars)
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

func (m *Meal) validatetGetAccessLevel(claims *models.Claims, vars *messages.RouteVars) (*uuid.UUID, *wrapper.HandlerError) {
	if vars == nil || vars.UserID == nil {
		if claims.AccessLevel == models.Admin {
			return nil, nil
		}
		return &claims.UserID, nil
	}

	if *vars.UserID == "me" || *vars.UserID == claims.UserID.String() {
		return &claims.UserID, nil
	}

	userID, err := uuid.Parse(*vars.UserID)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusBadRequest}
	}

	if claims.AccessLevel != models.Admin {
		return nil, &wrapper.HandlerError{Err: fmt.Errorf("cannot access other user records"), StatusCode: http.StatusUnauthorized}
	}

	return &userID, nil
}
