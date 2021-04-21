package models

import (
	"time"

	"github.com/google/uuid"
)

//User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID       uuid.UUID `json:"_" pg:"id, pk"`
	UserName string    `pg:"user_name,notnull"`
	Password string    `pg:"password,notnull"`
	Token    string    `pg:"access_token"`
	//	RefreshToken string

	CreatedAt time.Time `pg:"created_at,notnull"`
	UpdatedAt time.Time `pg:"updated_at,notnull"`
}
