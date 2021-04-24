package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/guilhermeCoutinho/api-studies/models"
)

func (u *Usecase) CreateMeal(ctx context.Context, req *messages.CreateMealRequest) error {
	claims, err := ClaimsFromCtx(ctx)
	if err != nil {
		return err
	}
	userID := claims.UserID

	meal, err := u.mealFromRequest(userID, req)
	if err != nil {
		return err
	}

	err = u.validateNewMealEntry(meal)
	if err != nil {
		return err
	}

	err = u.tryEnrichFromCaloriesAPI(meal)
	if err != nil {
		return err
	}

	return u.dal.Meal.UpsertMeal(ctx, meal)
}

func (u *Usecase) mealFromRequest(userID uuid.UUID, req *messages.CreateMealRequest) (*models.Meal, error) {
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

func (u *Usecase) validateNewMealEntry(meal *models.Meal) error {
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

func (u *Usecase) tryEnrichFromCaloriesAPI(meal *models.Meal) error {
	if meal.Calories == 0 {
		fetchedCalories, err := u.fetchCaloriesFromProvider(meal.Meal)
		if err != nil {
			return err
		}

		meal.Calories = fetchedCalories
	}
	return nil
}

func (u *Usecase) fetchCaloriesFromProvider(text string) (int, error) {
	return 99, nil
}
