package dal

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/models"
	"github.com/spf13/viper"
)

type UserDAL interface {
	UpsertUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, userName string, options *QueryOptions) (*models.User, error)
	GetUsers(ctx context.Context, userID *uuid.UUID, options *QueryOptions) ([]*models.User, error)
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
	_, err := u.db.Model(user).OnConflict("(id) DO UPDATE").Insert()
	return err
}

func (u *User) GetUser(
	ctx context.Context,
	userName string,
	options *QueryOptions,
) (*models.User, error) {
	user := &models.User{}

	partialQuery := u.db.Model(user).Where("user_name = ?", userName)
	partialQuery, err := addQueryOptions(partialQuery, options)
	if err != nil {
		return nil, err
	}

	err = partialQuery.Select()
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) GetUsers(
	ctx context.Context,
	userID *uuid.UUID,
	options *QueryOptions,
) ([]*models.User, error) {
	var users []*models.User

	partialQuery := u.db.Model(&users)
	if userID != nil {
		partialQuery = partialQuery.Where("id = ?", *userID)
	}

	partialQuery, err := addQueryOptions(partialQuery, options)
	if err != nil {
		return nil, err
	}

	err = partialQuery.Select()
	if err != nil {
		return nil, err
	}

	return users, err
}
