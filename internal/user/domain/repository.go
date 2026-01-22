package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserCreator interface {
	CreateUser(ctx context.Context, user *User) error
}

type UserByExternalIDGetter interface {
	GetUserByExternalID(ctx context.Context, externalID string) (*User, error)
}

type IDGenerator interface {
	Generate() (uuid.UUID, error)
}
