package models

import (
	"time"

	"github.com/google/uuid"
)

//User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID           uuid.UUID `json:"-" pg:"id, pk"`
	UserName     string    `json:"userName" pg:"user_name,notnull"`
	Password     string    `json:"-" pg:"password,notnull"`
	CalorieLimit int       `json:"caloriesLimit" pg:"calorie_limit, notnull"`

	CreatedAt time.Time `json:"-" pg:"created_at,notnull"`
	UpdatedAt time.Time `json:"-" pg:"updated_at,notnull"`
}
