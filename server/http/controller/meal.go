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
	"github.com/spf13/viper"
)

type Meal struct {
	dal    *dal.DAL
	config *viper.Viper
}

func NewMeal(
	dal *dal.DAL,
	config *viper.Viper,
) *Meal {
	return &Meal{
		dal:    dal,
		config: config,
	}
}

type Vars struct {
	UserID string `json:"userID"`
}

func (m *Meal) Post(ctx context.Context, args *messages.CreateMealRequest, vars *Vars) (*messages.BaseResponse, error) {
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

	meal, err := m.mealFromRequest(userID, args)
	if err != nil {
		return nil, err
	}

	err = m.validateNewMealEntry(meal)
	if err != nil {
		return nil, err
	}

	err = m.tryEnrichFromCaloriesAPI(meal)
	if err != nil {
		return nil, err
	}

	err = m.dal.Meal.UpsertMeal(ctx, meal)
	if err != nil {
		return nil, err
	}

	return &messages.BaseResponse{Code: http.StatusOK}, nil
}

func (m *Meal) mealFromRequest(userID uuid.UUID, req *messages.CreateMealRequest) (*models.Meal, error) {
	mealDate, err := time.Parse("2006-Jan-01", req.Date)
	if err != nil {
		return nil, err
	}

	mealTime, err := time.ParseDuration(req.Time)
	if err != nil {
		return nil, err
	}

	meal := &models.Meal{
		ID:       uuid.New(),
		UserID:   userID,
		Meal:     req.Meal,
		Calories: req.Calories,
		Date:     mealDate.Add(mealTime),

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

func (m *Meal) tryEnrichFromCaloriesAPI(meal *models.Meal) error {
	if meal.Calories == 0 {
		fetchedCalories, err := m.fetchCaloriesFromProvider(meal.Meal)
		if err != nil {
			return err
		}

		meal.Calories = fetchedCalories
	}
	return nil
}

func (m *Meal) fetchCaloriesFromProvider(text string) (int, error) {
	return 99, nil
}
