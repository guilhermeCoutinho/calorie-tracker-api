package meal

import (
	"context"
	"net/http"
	"time"

	"github.com/guilhermeCoutinho/api-studies/controller/contextextensions"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
)

func (m *Meal) Put(ctx context.Context, args *messages.UpdateMealRequest, vars *messages.RouteVars) (*messages.CreateMealResponse, *wrapper.HandlerError) {
	claims, err := contextextensions.ClaimsFromCtx(ctx)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	err = m.validatePostAccessLevel(claims, args.CreateMealPayload)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusUnauthorized}
	}

	meal, err := m.dal.Meal.GetMeal(ctx, args.ID, args.UserID)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusNotFound}
	}

	err = m.tryEnrichFromCaloriesAPI(args.CreateMealPayload)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	err = m.upadteMeal(meal, args.CreateMealPayload)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusBadRequest}
	}

	err = m.validateNewMealEntry(meal)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusBadRequest}
	}

	err = m.dal.Meal.UpsertMeal(ctx, meal)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	return &messages.CreateMealResponse{
		Meals: meal,
	}, nil
}

func (m *Meal) upadteMeal(meal *models.Meal, req *messages.CreateMealPayload) error {
	var err error
	if req.Calories != nil {
		meal.Calories = *req.Calories
	}

	if req.Date != nil {
		meal.Date, err = time.Parse("2006-Jan-02", *req.Date)
		if err != nil {
			return err
		}
	}

	if req.Meal != nil {
		meal.Meal = *req.Meal
	}

	if req.Time != nil {
		mealTime, err := time.ParseDuration(*req.Time)
		if err != nil {
			return err
		}
		meal.TimeSeconds = int64(mealTime.Seconds())
	}

	if req.UserID != nil {
		meal.UserID = *req.UserID
	}
	return nil
}
