package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	core "github.com/jacob-ebey/graphql-core"
)

const DefaultJwtHeader = "Bearer"

type JwtClaims struct {
	jwt.StandardClaims
	Claims
}

type JwtAuthHook struct {
	JwtHeader string
	JwtSecret []byte
}

func (hook *JwtAuthHook) GenerateToken(ctx context.Context, claims Claims, expirationMinutes int64) (string, error) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Duration(expirationMinutes) * time.Minute)

	t := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), JwtClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issuedAt.Unix(),
			ExpiresAt: expiresAt.Unix(),
		},
		Claims: Claims{
			ID:    claims.ID,
			Email: claims.Email,
			Role:  claims.Role,
		},
	})

	return t.SignedString(hook.JwtSecret)
}

func (hook *JwtAuthHook) PreExecute(ctx context.Context, req core.GraphQLRequest) context.Context {
	authorization, _ := ctx.Value("authorization").(string)

	var claims *Claims = nil

	if token := hook.getToken(authorization); token != "" {
		claims = hook.GetClimasFromToken(token)
	}

	ctx = context.WithValue(ctx, "claims", claims)

	return context.WithValue(ctx, "auth", hook)
}

func (hook *JwtAuthHook) getToken(str string) string {
	split := strings.SplitN(str, " ", 2)

	header := hook.JwtHeader
	if header == "" {
		header = DefaultJwtHeader
	}

	if len(split) == 2 && split[0] == header {
		return split[1]
	}

	return ""
}

func (hook *JwtAuthHook) GetClimasFromToken(token string) *Claims {
	parsed, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return hook.JwtSecret, nil
	})

	if err != nil {
		return nil
	}

	if claims, ok := parsed.Claims.(*JwtClaims); ok && parsed.Valid && claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return &claims.Claims
	}

	return nil
}
