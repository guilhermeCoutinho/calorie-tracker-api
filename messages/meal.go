package messages

import (
	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/models"
)

type CreateMealPayload struct {
	UserID   *uuid.UUID
	Meal     string `json:"meal"`
	Calories *int   `json:"calories"`
	Date     string `json:"date"`
	Time     string `json:"time"`
}

type CreateMealResponse struct {
	BaseResponse
	Meals *models.Meal `json:"meal"`
}

type RouteVars struct {
	UserID *string `json:"userID"`
}

type GetMealsResponse struct {
	BaseResponse
	Meals []*models.MealWithLimit `json:"meals"`
}
