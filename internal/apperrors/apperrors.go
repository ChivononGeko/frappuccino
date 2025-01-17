package apperrors

import "errors"

var (
	ErrInvalidInput     = errors.New("invalid input: missing required field")
	ErrExistConflict    = errors.New("already exist")
	ErrNotExistConflict = errors.New("doesn't exist")
	ErrOrderClosed      = errors.New("the order is already closed")
)
