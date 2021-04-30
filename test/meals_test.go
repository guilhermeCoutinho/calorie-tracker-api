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
	t.Parallel()
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

	doRequest(t, http.MethodPost, "/meals", &token, createMealRequest, &map[string]interface{}{})
	fromDBResponse := &messages.GetMealsResponse{}
	doRequest(t, http.MethodGet, "/meals", &token, struct{}{}, fromDBResponse)

	assert.Equal(t, 100, fromDBResponse.Meals[0].Calories)
	assert.Equal(t, "hamburguer", fromDBResponse.Meals[0].Meal.Meal)
}

func TestCreateMeal(t *testing.T) {
	t.Parallel()
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

	doRequest(t, http.MethodPost, "/meals", &token, createMealRequest, &map[string]interface{}{})
	fromDBResponse := &messages.GetMealsResponse{}
	doRequest(t, http.MethodGet, "/meals", &token, struct{}{}, fromDBResponse)

	assert.Equal(t, 100, fromDBResponse.Meals[0].Calories)
	assert.Equal(t, "hamburguer", fromDBResponse.Meals[0].Meal.Meal)
}

func TestCreateMealWithCalorieProvider(t *testing.T) {
	t.Parallel()
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

	fromDBResponse := &messages.GetMealsResponse{}
	doRequest(t, http.MethodGet, "/meals", &token, struct{}{}, fromDBResponse)

	assert.NotZero(t, fromDBResponse.Meals[0].Calories)
}

func TestCreateMealFail(t *testing.T) {
	t.Parallel()
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

	statusCode := doRequest(t, http.MethodPost, "/users/me/meals", &token, createMealRequest, &map[string]interface{}{})
	assert.Equal(t, http.StatusBadRequest, statusCode)
}

func TestMealsBelowLimit(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
