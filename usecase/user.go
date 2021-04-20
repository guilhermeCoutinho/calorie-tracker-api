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
