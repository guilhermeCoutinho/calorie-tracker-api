package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/models"
)

func (u *Usecase) CreateUser(userName, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		ID:       uuid.New(),
		UserName: userName,
		Password: hashedPassword,
	}

	return u.dal.UpsertUser(context.Background(), user)
}

func (u *Usecase) UserLogin(userName, password string) (*models.User, error) {
	ctx := context.Background()
	user, err := u.dal.GetUser(ctx, userName)
	if err != nil {
		return nil, err
	}

	err = verifyPassword(user.Password, password)
	if err != nil {
		return nil, err
	}

	user.Token = uuid.New().String()
	err = u.dal.UpsertUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
