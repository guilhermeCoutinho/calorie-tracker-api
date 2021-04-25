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

type GetMealsVars struct {
	UserID string `json:"userID"`
	Filter string `json:"filter"` //(date eq '2016-05-01') AND ((number_of_calories gt 20) OR (number_of_calories lt 10))
}

type GetMealsResponse struct {
	BaseResponse
	Meals []*models.MealWithLimit `json:"meals"`
}
