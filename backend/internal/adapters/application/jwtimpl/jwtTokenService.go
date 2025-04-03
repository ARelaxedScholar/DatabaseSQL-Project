package jwtimpl

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sql-project-backend/internal/ports"
)

type JwtTokenService struct {
	secretKey     []byte
	tokenDuration time.Duration
}

func NewJwtTokenService(secretKey string, tokenDuration time.Duration) ports.TokenService {
	return &JwtTokenService{
		secretKey:     []byte(secretKey),
		tokenDuration: tokenDuration,
	}
}

type jwtCustomClaims struct {
	UserID int    `json:"userId"`
	Role   string `json:"role"` // e.g., "client", "employee", "admin"
	jwt.RegisteredClaims
}

func (s *JwtTokenService) GenerateTokenWithDuration(userID int, role string, duration time.Duration) (string, error) {
	claims := &jwtCustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "SunflowerBooking",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JwtTokenService) ValidateToken(tokenString string) (int, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})
	if err != nil {
		return 0, "", err
	}
	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		return claims.UserID, claims.Role, nil
	}
	return 0, "", errors.New("invalid token")
}
