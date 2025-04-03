package auth

import (
	"strings"

	"github.com/sql-project-backend/internal/ports"
)

// contextKey is a custom type to avoid context key collisions.
type contextKey string

// ContextUserID is the key used to store the authenticated user's ID in the context.
const ContextUserID contextKey = "userID"

// ValidateAuthToken centralizes token validation.
// It strips a "Bearer " prefix if present, validates the token using the provided tokenService,
// and returns the user ID (or an error if validation fails).
func ValidateAuthToken(tokenString string, tokenService ports.TokenService) (int, string, error) {
	// Remove any Bearer prefix (case insensitive)
	if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
		tokenString = tokenString[len("bearer "):]
	}
	// Validate the token using your token service.
	return tokenService.ValidateToken(tokenString)
}
