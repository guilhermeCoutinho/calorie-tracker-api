package meal

import (
	"github.com/guilhermeCoutinho/api-studies/dal"
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
