package utils

import "github.com/google/uuid"

type UUID7Generator struct{}

func NewUUID7Generator() *UUID7Generator {
	return &UUID7Generator{}
}

func (g *UUID7Generator) GenID() (uuid.UUID, error) {
	return uuid.NewV7()
}
