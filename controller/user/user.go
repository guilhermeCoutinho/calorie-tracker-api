package user

import (
	"github.com/guilhermeCoutinho/api-studies/dal"
	"github.com/spf13/viper"
)

type User struct {
	dal    *dal.DAL
	config *viper.Viper
}

func NewUser(
	dal *dal.DAL,
	config *viper.Viper,
) *User {
	return &User{
		dal:    dal,
		config: config,
	}
}
