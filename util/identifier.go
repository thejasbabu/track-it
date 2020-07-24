package util

import "github.com/google/uuid"

// Identifier to generate unique id
type Identifier interface {
	Generate() string
}

// UUIDIdentifier uses google uuid to generate a random (Version 4) UUID
type UUIDIdentifier struct{}

// Generate returns a google uuid
func (UUIDIdentifier) Generate() string {
	return uuid.New().String()
}
