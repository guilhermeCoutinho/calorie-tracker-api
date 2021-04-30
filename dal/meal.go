package dal

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/spf13/viper"
)

type MealDAL interface {
	UpsertMeal(ctx context.Context, user *models.Meal) error
	GetMeals(ctx context.Context, userID *uuid.UUID, options *QueryOptions) ([]*models.MealWithLimit, error)
	GetMeal(ctx context.Context, id uuid.UUID, userID *uuid.UUID) (*models.Meal, error)
	DeleteMeal(ctx context.Context, id uuid.UUID, userID *uuid.UUID) error
}

type Meal struct {
	config *viper.Viper
	db     *pg.DB
}

func NewMeal(
	config *viper.Viper,
	db *pg.DB,
) *Meal {
	return &Meal{
		config: config,
		db:     db,
	}
}

func (u *Meal) UpsertMeal(ctx context.Context, meal *models.Meal) error {
	meal.UpdatedAt = time.Now()
	_, err := u.db.Model(meal).OnConflict("(id) DO UPDATE").Insert()
	return err
}

func (u *Meal) GetMeals(ctx context.Context, userID *uuid.UUID, options *QueryOptions) ([]*models.MealWithLimit, error) {
	var meals []*models.MealWithLimit

	partialQuery := u.db.Model(&meals)
	if userID != nil {
		partialQuery = partialQuery.Where("user_id = ?", *userID)
	}

	partialQuery, err := addQueryOptions(partialQuery, options)
	if err != nil {
		return nil, err
	}

	err = partialQuery.Select()
	if err != nil {
		return nil, err
	}

	return meals, err
}

func (u *Meal) GetMeal(ctx context.Context, id uuid.UUID, userID *uuid.UUID) (*models.Meal, error) {
	var meal *models.Meal

	partialQuery := u.db.Model(&meal).Where("id = ?", id)
	if userID != nil {
		partialQuery = partialQuery.Where("user_id = ?", *userID)
	}

	err := partialQuery.Select()
	if err != nil {
		return nil, err
	}

	return meal, err
}

func (u *Meal) DeleteMeal(ctx context.Context, id uuid.UUID, userID *uuid.UUID) error {
	meal := models.Meal{
		ID: id,
	}

	partialQuery := u.db.Model(&meal).Where("id = ?", id)
	if userID != nil {
		partialQuery = partialQuery.Where("user_id = ?", *userID)
	}

	result, err := partialQuery.Delete()
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no rows")
	}
	return err
}
