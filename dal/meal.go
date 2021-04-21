package dal

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/spf13/viper"
)

type MealDAL interface {
	UpsertMeal(ctx context.Context, user *models.Meal) error
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
