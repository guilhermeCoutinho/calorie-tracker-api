package messages

import "github.com/guilhermeCoutinho/api-studies/models"

type CreateMealPayload struct {
	Meal     string `json:"meal"`
	Calories *int   `json:"calories"`
	Date     string `json:"date"`
	Time     string `json:"time"`
}

type CreateMealVars struct {
	UserID string `json:"userID"`
}

type CreateMealResponse struct {
	BaseResponse
	Meals *models.Meal `json:"meal"`
}

type GetMealsVars struct {
	UserID     string `json:"userID"`
	Pagination string `json:"pagination"`
}

type GetMealsResponse struct {
	BaseResponse
	Meals []*models.MealWithLimit `json:"meals"`
}
