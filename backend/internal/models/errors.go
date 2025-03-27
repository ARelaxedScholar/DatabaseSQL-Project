package models

import "errors"

// Standard repository errors
var (
	ErrNotFound = errors.New("Requested record not found.")
	ErrDuplicateEntry = errors.New("Database constraint violation: duplicate entry.")
)