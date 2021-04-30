package models

import (
	"time"

	"github.com/google/uuid"
)

type Meal struct {
	ID          uuid.UUID `json:"id" pg:"id, pk"`
	UserID      uuid.UUID `json:"user_id" pg:"user_id"`
	Meal        string    `json:"meal" pg:"meal"`
	Calories    int       `json:"calories" pg:"calories"`
	Date        time.Time `json:"date" pg:"date"`
	TimeSeconds int64     `json:"time_seconds" pg:"time_seconds"`

	CreatedAt time.Time `json:"-" pg:"created_at,notnull"`
	UpdatedAt time.Time `json:"-" pg:"updated_at,notnull"`
}

type MealWithLimit struct {
	tableName struct{} `pg:"select:meals_with_limit"`
	*Meal
	AboveCaloriesLimit  bool `json:"above_limit" pg:"above_limit"`
	TotalCaloriesForDay int  `json:"total_calories_for_day" pg:"total_calories_for_day"`
}
