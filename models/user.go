package models

import (
	"time"

	"github.com/google/uuid"
)

//User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID           uuid.UUID   `json:"id" pg:"id, pk"`
	UserName     string      `json:"user_name" pg:"user_name,notnull"`
	Password     string      `json:"-" pg:"password,notnull"`
	CalorieLimit int         `json:"calorie_limit" pg:"calorie_limit" sql:",notnull"`
	AccessLevel  AccessLevel `json:"access_level" pg:"access_level, notnull" sql:",notnull"`

	CreatedAt time.Time `json:"-" pg:"created_at,notnull"`
	UpdatedAt time.Time `json:"-" pg:"updated_at,notnull"`
}

type AccessLevel int

const (
	Admin       AccessLevel = 0
	Manager     AccessLevel = 1
	RegulerUser AccessLevel = 2
)

type Claims struct {
	UserID      uuid.UUID
	AccessLevel AccessLevel
}
