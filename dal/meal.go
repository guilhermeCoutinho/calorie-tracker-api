package dal

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/spf13/viper"
)

type MealDAL interface {
	UpsertMeal(ctx context.Context, user *models.Meal) error
	GetMeals(ctx context.Context, userID uuid.UUID, query string) ([]*models.MealWithLimit, error)
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
	query := u.db.Model(meal).OnConflict("(id) DO UPDATE")
	err := upsertAllFields(query, meal)
	return err
}

func (u *Meal) GetMeals(ctx context.Context, userID uuid.UUID, query string) ([]*models.MealWithLimit, error) {
	var meals []*models.MealWithLimit
	err := u.db.Model(&meals).Where("user_id=?", userID.String()).Select()
	return meals, err
}
