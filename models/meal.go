package models

import (
	"time"

	"github.com/google/uuid"
)

type Meal struct {
	ID       uuid.UUID `json:"-" pg:"id, pk"`
	UserID   uuid.UUID `pg:"user_id"`
	Meal     string    `pg:"meal"`
	Calories int       `pg:"calories"`
	Date     time.Time `pg:"date"`

	CreatedAt time.Time `json:"-" pg:"created_at,notnull"`
	UpdatedAt time.Time `json:"-" pg:"updated_at,notnull"`
}

type MealWithLimit struct {
	tableName struct{} `pg:"select:meals_with_limit"`
	*Meal
	AboveCaloriesLimit bool `pg:"above_limit"`
}
