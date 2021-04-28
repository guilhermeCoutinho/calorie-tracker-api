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

// EXP	: (EXP)
// 		| EXP (OR/AND) EXP
//		| columnName op value
func (m *Meal) Post(ctx context.Context, args *messages.CreateMealPayload, vars *messages.CreateMealVars) (*messages.CreateMealResponse, error) {
	claims, err := ClaimsFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	userID := uuid.Nil
	if vars.UserID == "me" {
		userID = claims.UserID
	} else {
		userID, err = uuid.Parse(vars.UserID)
		if err != nil {
			return nil, err
		}
	}

	err = m.tryEnrichFromCaloriesAPI(args)
	if err != nil {
		return nil, err
	}

	meal, err := m.mealFromRequest(userID, args)
	if err != nil {
		return nil, err
	}

	err = m.validateNewMealEntry(meal)
	if err != nil {
		return nil, err
	}

	err = m.dal.Meal.UpsertMeal(ctx, meal)
	if err != nil {
		return nil, err
	}

	return &messages.CreateMealResponse{
		BaseResponse: messages.BaseResponse{Code: http.StatusOK},
		Meals:        meal,
	}, nil
}

func (m *Meal) Get(ctx context.Context, args *messages.GetMealsResponse, vars *messages.GetMealsVars) (*messages.GetMealsResponse, error) {
	claims, err := ClaimsFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	userID := uuid.Nil
	if vars.UserID == "me" {
		userID = claims.UserID
	} else {
		userID, err = uuid.Parse(vars.UserID)
		if err != nil {
			return nil, err
		}
	}

	meals, err := m.dal.Meal.GetMeals(ctx, userID, getQueryOptions(ctx))
	if err != nil {
		return nil, err
	}

	return &messages.GetMealsResponse{
		BaseResponse: messages.BaseResponse{
			Code: http.StatusOK,
		},
		Meals: meals,
	}, nil
}

func (m *Meal) mealFromRequest(userID uuid.UUID, req *messages.CreateMealPayload) (*models.Meal, error) {
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
		UserID:      userID,
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
