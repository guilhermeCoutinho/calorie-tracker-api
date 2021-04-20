package usecase

import (
	"github.com/guilhermeCoutinho/api-studies/dal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Usecase struct {
	dal    dal.UserDAL
	config *viper.Viper
	logger logrus.FieldLogger
}

func NewUsecase(
	config *viper.Viper,
	dal dal.UserDAL,
	logger logrus.FieldLogger,
) *Usecase {
	return &Usecase{
		config: config,
		dal:    dal,
		logger: logger,
	}
}
