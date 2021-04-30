package messages

import (
	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/models"
)

type CreateUserRequest struct {
	Username     string `json:"user_name"`
	Password     string `json:"password"`
	CalorieLimit int    `json:"calorie_limit"`
}

type GetUsersResponse struct {
	Users []*models.User `json:"users"`
}

type UpdateUserRequest struct {
	ID           *uuid.UUID
	Username     *string             `json:"user_name"`
	CalorieLimit *int                `json:"calorie_limit"`
	AccessLevel  *models.AccessLevel `json:"access_level"`
	Password     *string             `json:"password"`
}
