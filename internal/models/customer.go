package models

import (
	"encoding/json"
	"frappuchino/internal/apperrors"
)

type Customer struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Preferences json.RawMessage `json:"preferences"`
}

func NewCustomer(name string, instructions json.RawMessage) (*Customer, error) {
	if name == "" {
		return nil, apperrors.ErrInvalidInput
	}

	if instructions == nil {
		instructions = json.RawMessage(`{"preferences": true}`)
	}

	email := fromNameToID(name) + "@gmail.com"

	return &Customer{
		Name:        name,
		Email:       email,
		Preferences: instructions,
	}, nil
}
