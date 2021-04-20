package dal

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/spf13/viper"
)

type UserDAL interface {
	UpsertUser(ctx context.Context, user *models.User) error
}

type User struct {
	config *viper.Viper
	db     *pg.DB
}

func NewUser(
	config *viper.Viper,
	db *pg.DB,
) *User {
	return &User{
		config: config,
		db:     db,
	}
}

func (u *User) UpsertUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	query := u.db.Model(user).OnConflict("(id) DO UPDATE")
	err := upsertAllFields(query, user)
	return err
}
