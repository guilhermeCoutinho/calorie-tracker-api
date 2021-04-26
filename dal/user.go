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

type UserDAL interface {
	UpsertUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, userName string, options *QueryOptions) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID, options *QueryOptions) (*models.User, error)
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

func (u *User) GetUser(
	ctx context.Context,
	userName string,
	options *QueryOptions,
) (*models.User, error) {
	return u.getUser(ctx, "user_name", userName, options)
}

func (u *User) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
	options *QueryOptions,
) (*models.User, error) {
	return u.getUser(ctx, "id", userID, options)
}

func (u *User) getUser(
	ctx context.Context,
	column string,
	value interface{},
	options *QueryOptions,
) (*models.User, error) {
	user := &models.User{}
	condition := fmt.Sprintf("%s=?", column)

	partialQuery := u.db.Model(user).Where(condition, value)
	err := addQueryOptions(partialQuery, options).Select()
	if err != nil {
		return nil, err
	}
	return user, nil
}
