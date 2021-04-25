package messages

import "github.com/guilhermeCoutinho/api-studies/models"

type CreateUserRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	CalorieLimit int    `json:"calorieLimit"`
}

type GetUsersResponse struct {
	BaseResponse
	Users *models.User `json:"users"`
}

type UpdateUserRequest struct {
	Username     *string `json:"userName"`
	CalorieLimit *int    `json:"caloriesLimit"`
}
