package models

import (
	"time"

	"github.com/google/uuid"
)

//User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID       uuid.UUID `json:"_"`
	UserName string
	Password string
	//	Token         string
	//	RefreshToken string

	CreatedAt time.Time
	UpdatedAt time.Time
}
