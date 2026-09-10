package grpc

import (
	"errors"
	"fmt"
	"time"

	jwtpkg "github.com/golang-jwt/jwt/v5"
)

func generateJWT(secret, userID, email string) (string, error) {
	if secret == "" {
		return "", errors.New("jwt secret is required")
	}
	if userID == "" {
		return "", errors.New("user id is required")
	}
	if email == "" {
		return "", errors.New("email is required")
	}

	claims := jwtpkg.MapClaims{
		"sub":   userID,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(30 * time.Minute).Unix(),
		"type":  "access",
	}

	token := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func validateJWT(secret, tokenString string) (string, string, error) {
	if secret == "" {
		return "", "", errors.New("jwt secret is required")
	}

	parsed, err := jwtpkg.Parse(tokenString, func(token *jwtpkg.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtpkg.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", "", err
	}

	claims, ok := parsed.Claims.(jwtpkg.MapClaims)
	if !ok || !parsed.Valid {
		return "", "", errors.New("invalid token")
	}

	userID, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	return userID, email, nil
}
