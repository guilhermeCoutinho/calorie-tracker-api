package dal

import (
	"github.com/go-pg/pg/v10"
	"github.com/spf13/viper"
)

type DAL struct {
	Meal MealDAL
	User UserDAL
}

func NewDAL(
	config *viper.Viper,
	db *pg.DB,
) *DAL {
	return &DAL{
		Meal: NewMeal(config, db),
		User: NewUser(config, db),
	}
}
