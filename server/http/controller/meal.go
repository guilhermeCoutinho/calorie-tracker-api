package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/dal"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/guilhermeCoutinho/api-studies/server/http/wrapper"
	"github.com/guilhermeCoutinho/api-studies/services/calorieprovider"
	"github.com/spf13/viper"
)

type Meal struct {
	dal            *dal.DAL
	config         *viper.Viper
	calorieService calorieprovider.Provider
}

func NewMeal(
	dal *dal.DAL,
	config *viper.Viper,
	calorieService calorieprovider.Provider,
) *Meal {
	return &Meal{
		dal:            dal,
		config:         config,
		calorieService: calorieService,
	}
}

func (m *Meal) Post(ctx context.Context, args *messages.CreateMealPayload, vars *messages.RouteVars) (*messages.CreateMealResponse, *wrapper.HandlerError) {
	claims, err := ClaimsFromCtx(ctx)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	err = m.validatePostAccessLevel(claims, args)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusUnauthorized}
	}

	err = m.tryEnrichFromCaloriesAPI(args)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	meal, err := m.mealFromRequest(args)
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
		BaseResponse: messages.BaseResponse{Code: http.StatusOK},
		Meals:        meal,
	}, nil
}

func (m *Meal) validatePostAccessLevel(claims *Claims, args *messages.CreateMealPayload) error {
	if args.UserID == nil {
		args.UserID = &claims.UserID
	}

	if claims.AccessLevel == models.Admin {
		return nil
	}

	if *args.UserID != claims.UserID {
		return fmt.Errorf("wrong access level")
	}

	return nil
}

func (m *Meal) Get(ctx context.Context, args *messages.GetMealsResponse, vars *messages.RouteVars) (*messages.GetMealsResponse, *wrapper.HandlerError) {
	claims, err := ClaimsFromCtx(ctx)
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	userID, wrapperErr := m.validatetGetAccessLevel(claims, vars)
	fmt.Println("userID to filter is ", userID)
	if wrapperErr != nil {
		return nil, wrapperErr
	}

	meals, err := m.dal.Meal.GetMeals(ctx, userID, getQueryOptions(ctx))
	if err != nil {
		return nil, &wrapper.HandlerError{Err: err, StatusCode: http.StatusInternalServerError}
	}

	return &messages.GetMealsResponse{
		BaseResponse: messages.BaseResponse{
			Code: http.StatusOK,
		},
		Meals: meals,
	}, nil
}

func (m *Meal) validatetGetAccessLevel(claims *Claims, vars *messages.RouteVars) (*uuid.UUID, *wrapper.HandlerError) {
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

func (m *Meal) mealFromRequest(req *messages.CreateMealPayload) (*models.Meal, error) {
	mealDate, err := time.Parse("2006-Jan-02", req.Date)
	if err != nil {
		return nil, err
	}

	mealTime, err := time.ParseDuration(req.Time)
	if err != nil {
		return nil, err
	}

	meal := &models.Meal{
		ID:          uuid.New(),
		UserID:      *req.UserID,
		Meal:        req.Meal,
		Calories:    *req.Calories,
		Date:        mealDate,
		TimeSeconds: int64(mealTime.Seconds()),

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return meal, nil
}

func (m *Meal) validateNewMealEntry(meal *models.Meal) error {
	if meal.Date.After(time.Now().UTC()) {
		return fmt.Errorf("cannot meals in the future")
	}

	if meal.Calories < 0 {
		return fmt.Errorf("calories cannot be negative")
	}

	if meal.Meal == "" {
		return fmt.Errorf("meal cannot have empty text")
	}

	return nil
}

func (m *Meal) tryEnrichFromCaloriesAPI(meal *messages.CreateMealPayload) error {
	if meal.Calories == nil {
		fetchedCalories, err := m.calorieService.GetCalories(meal.Meal)
		if err != nil {
			return err
		}

		meal.Calories = &fetchedCalories
	}
	return nil
}
