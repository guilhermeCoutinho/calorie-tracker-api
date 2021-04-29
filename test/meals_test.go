package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/guilhermeCoutinho/api-studies/messages"
	"github.com/stretchr/testify/assert"
)

const URL = "http://localhost:8080"

func TestCreateMealSuccess(t *testing.T) {
	drop := GetPG(t)
	drop()

	token := getAuthenticatedUser(t)

	calories := 100
	createMealRequest := &messages.CreateMealPayload{
		Meal:     "hamburguer",
		Calories: &calories,
		Date:     "2021-Jan-01",
		Time:     "2h00m",
	}

	createMealResponse := &messages.CreateMealResponse{}
	doRequest(t, http.MethodPost, "/users/me/meals", &token, createMealRequest, createMealResponse)

	assert.NotNil(t, createMealResponse.Meals)
	assert.Equal(t, calories, createMealResponse.Meals.Calories)
	assert.Equal(t, "hamburguer", createMealResponse.Meals.Meal)
	assert.Equal(t, "2021-01-01 00:00:00 +0000 UTC", createMealResponse.Meals.Date.String())
}

func TestCreateMealWithCalorieProvider(t *testing.T) {
	drop := GetPG(t)
	drop()

	token := getAuthenticatedUser(t)

	createMealRequest := &messages.CreateMealPayload{
		Meal: "hamburguer",
		Date: "2021-Jan-01",
		Time: "2h00m",
	}

	createMealResponse := &messages.CreateMealResponse{}
	doRequest(t, http.MethodPost, "/users/me/meals", &token, createMealRequest, createMealResponse)

	assert.NotZero(t, createMealResponse.Meals.Calories)
}

func TestCreateMealFail(t *testing.T) {
	drop := GetPG(t)
	drop()

	token := getAuthenticatedUser(t)

	calories := -10
	createMealRequest := &messages.CreateMealPayload{
		Meal:     "hamburguer",
		Date:     "2021-Jan-01",
		Time:     "2h00m",
		Calories: &calories,
	}

	createMealResponse := &messages.CreateMealResponse{}
	statusCode := doRequest(t, http.MethodPost, "/users/me/meals", &token, createMealRequest, createMealResponse)

	assert.Equal(t, http.StatusInternalServerError, statusCode)
}

func TestMealsBelowLimit(t *testing.T) {
	drop := GetPG(t)
	drop()

	token := getAuthenticatedUser(t)

	calories := 99
	for i := 0; i < 5; i++ {
		createMealRequest := &messages.CreateMealPayload{
			Meal:     "hamburguer",
			Calories: &calories,
			Date:     fmt.Sprintf("2021-Jan-0%d", i+1),
			Time:     "12h",
		}

		createMealResponse := &messages.CreateMealResponse{}
		doRequest(t, http.MethodPost, "/users/me/meals", &token, createMealRequest, createMealResponse)
	}

	getMealsResponse := &messages.GetMealsResponse{}
	doRequest(t, http.MethodGet, "/users/me/meals", &token, &struct{}{}, getMealsResponse)
	assert.NotNil(t, getMealsResponse)

	for _, meal := range getMealsResponse.Meals {
		assert.False(t, meal.AboveCaloriesLimit)
	}
}

func TestMealsAboveLimit(t *testing.T) {
	drop := GetPG(t)
	drop()

	token := getAuthenticatedUser(t)

	calories := 50
	for i := 0; i < 3; i++ {
		createMealRequest := &messages.CreateMealPayload{
			Meal:     "hamburguer",
			Calories: &calories,
			Date:     "2021-Jan-01",
			Time:     "12h",
		}

		createMealResponse := &messages.CreateMealResponse{}
		doRequest(t, http.MethodPost, "/users/me/meals", &token, createMealRequest, createMealResponse)
	}

	getMealsResponse := &messages.GetMealsResponse{}
	doRequest(t, http.MethodGet, "/users/me/meals", &token, &struct{}{}, getMealsResponse)
	assert.NotNil(t, getMealsResponse)

	for _, meal := range getMealsResponse.Meals {
		assert.True(t, meal.AboveCaloriesLimit)
		assert.Equal(t, 150, meal.TotalCaloriesForDay)
	}
}

func printJSON(v interface{}) {
	bytes, _ := json.MarshalIndent(v, "", " ")
	fmt.Println(string(bytes))
}
