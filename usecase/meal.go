package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/models"
)

func (u *Usecase) CreateMeal(ctx context.Context, userToken uuid.UUID, text string, calories int) error {
	shouldFetchFromCaloriesAPI := calories <= 0

	if shouldFetchFromCaloriesAPI {
		fetchedCalories, err := u.tryAndFetchCalories(text)
		if err != nil {
			return err
		}

		calories = fetchedCalories
	}

	meal := &models.Meal{
		ID:       uuid.New(),
		UserID:   userToken,
		Metadata: text,
		Calories: calories,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return u.dal.Meal.UpsertMeal(context.Background(), meal)
}

func (u *Usecase) tryAndFetchCalories(text string) (int, error) {
	return -1, nil
}
