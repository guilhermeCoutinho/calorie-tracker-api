package models

import (
	"time"

	"github.com/google/uuid"
)

//User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID       uuid.UUID `json:"_" sql:"id, pk"`
	UserName string    `sql:"user_name,notnull"`
	Password string    `sql:"password,notnull"`
	//	Token         string
	//	RefreshToken string

	CreatedAt time.Time `sql:"created_at,notnull"`
	UpdatedAt time.Time `sql:"updated_at,notnull"`
}
