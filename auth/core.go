package auth

import (
	"context"
)

type Claims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type TokenGenerator interface {
	GenerateToken(ctx context.Context, claims Claims, expirationMinutes int64) (string, error)
}
