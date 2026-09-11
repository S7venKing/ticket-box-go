package auth

import (
	"errors"
	"fmt"
	"time"

	jwtpkg "github.com/golang-jwt/jwt/v5"
)

const DefaultTokenTTL = 30 * time.Minute

type Config struct {
	SecretKey string
	TTL       time.Duration
}

func (c Config) tokenTTL() time.Duration {
	if c.TTL <= 0 {
		return DefaultTokenTTL
	}
	return c.TTL
}

func NewTokenService(secret string) *TokenService {
	return &TokenService{secret: secret}
}

type TokenService struct {
	secret string
}

func (s *TokenService) Generate(userID, email string) (string, error) {
	return s.GenerateWithRole(userID, email, "user")
}

func (s *TokenService) GenerateWithRole(userID, email, role string) (string, error) {
	if s.secret == "" {
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
		"role":  role,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(DefaultTokenTTL).Unix(),
		"type":  "access",
	}

	token := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *TokenService) Validate(tokenString string) (Claims, error) {
	if s.secret == "" {
		return Claims{}, errors.New("jwt secret is required")
	}

	token, err := jwtpkg.Parse(tokenString, func(token *jwtpkg.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtpkg.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return Claims{}, err
	}

	claims, ok := token.Claims.(jwtpkg.MapClaims)
	if !ok || !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	userID, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)
	tokenType, _ := claims["type"].(string)
	if userID == "" || email == "" {
		return Claims{}, errors.New("invalid token claims")
	}
	if tokenType != "access" {
		return Claims{}, errors.New("unexpected token type")
	}

	if role == "" {
		role = "user"
	}
	return Claims{UserID: userID, Email: email, Role: role}, nil
}

type Claims struct {
	UserID string
	Email  string
	Role   string
}
