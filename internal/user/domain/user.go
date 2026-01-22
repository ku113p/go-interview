package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID

	CreatedAt time.Time

	ExternalID string
}

func NewUser(id uuid.UUID, externalID string) *User {
	return &User{
		ID:         id,
		CreatedAt:  time.Now(),
		ExternalID: externalID,
	}
}
