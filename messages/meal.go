package messages

import (
	"github.com/google/uuid"
	"github.com/guilhermeCoutinho/api-studies/models"
)

type CreateMealPayload struct {
	UserID   *uuid.UUID
	Meal     *string `json:"meal"`
	Calories *int    `json:"calories"`
	Date     *string `json:"date"`
	Time     *string `json:"time"`
}

type UpdateMealRequest struct {
	ID uuid.UUID
	*CreateMealPayload
}

type CreateMealResponse struct {
	Meals *models.Meal `json:"meal"`
}

type RouteVars struct {
	UserID *string    `json:"userID"`
	MealID *uuid.UUID `json:"mealID"`
}

type GetMealsResponse struct {
	Meals []*models.MealWithLimit `json:"meals"`
}

type DeleteMealRequest struct {
	ID uuid.UUID `json:"id"`
}
