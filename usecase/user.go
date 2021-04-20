package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/models"
)

func CreateUser(userName, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		ID:       uuid.New(),
		UserName: userName,
		Password: hashedPassword,
	}

	fmt.Println("Create user called with %s %+v", hashedPassword, user)

	return nil
}
