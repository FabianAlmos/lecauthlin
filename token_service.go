package main

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"strings"
	"time"
)

type JwtCustomClaims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func GenerateToken(userID, lifetimeMinutes int, secret string) (string, error) {
	claims := &JwtCustomClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(lifetimeMinutes)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateBearerToken(authHeader string) (*JwtCustomClaims, error) {
	tokenString := BearerAuthHeader(authHeader)

	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("failed to parse token claims")
	}

	return claims, nil
}

func BearerAuthHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, "Bearer")
	if len(parts) != 2 {
		return ""
	}

	token := strings.TrimSpace(parts[1])
	if len(token) < 1 {
		return ""
	}

	return token
}
