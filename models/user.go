package models

import (
	"time"

	"github.com/google/uuid"
)

//User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID           uuid.UUID   `json:"-" pg:"id, pk" sortable:"true"`
	UserName     string      `json:"userName" pg:"user_name,notnull"`
	Password     string      `json:"-" pg:"password,notnull"`
	CalorieLimit int         `json:"caloriesLimit" pg:"calorie_limit, notnull"`
	AccessLevel  AccessLevel `json:"acessLevel" pg:"access_level"`

	CreatedAt time.Time `json:"-" pg:"created_at,notnull"`
	UpdatedAt time.Time `json:"-" pg:"updated_at,notnull"`
}

type AccessLevel int

const (
	Admin       AccessLevel = 0
	Manager     AccessLevel = 1
	RegulerUser AccessLevel = 2
)
