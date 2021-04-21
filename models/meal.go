package models

import (
	"time"

	"github.com/google/uuid"
)

type Meal struct {
	ID                 uuid.UUID `json:"_" pg:"id, pk"`
	UserID             uuid.UUID `pg:"user_id"`
	Metadata           string    `pg:"metadata"`
	Calories           int       `pg:"calories"`
	AboveCaloriesLimit bool      `pg:"more_than_limit"`

	CreatedAt time.Time `pg:"created_at,notnull"`
	UpdatedAt time.Time `pg:"updated_at,notnull"`
}
