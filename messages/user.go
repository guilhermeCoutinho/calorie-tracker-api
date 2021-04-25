package messages

type CreateUserRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	CalorieLimit int    `json:"calorieLimit"`
}

type UpdateUserRequest struct {
	Username     *string `json:"userName"`
	CalorieLimit *int    `json:"calorieLimit"`
}
