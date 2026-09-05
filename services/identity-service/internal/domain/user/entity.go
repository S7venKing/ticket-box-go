package user

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID

	Email        string
	PasswordHash string

	FullName string
	Phone    string

	IsActive bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(
	email string,
	passwordHash string,
	fullName string,
	phone string,
) (*User, error) {
	email = normalizeEmail(email)
	fullName = strings.TrimSpace(fullName)
	phone = strings.TrimSpace(phone)

	if email == "" {
		return nil, ErrEmailRequired
	}

	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	if passwordHash == "" {
		return nil, ErrPasswordHashRequired
	}

	if fullName == "" {
		return nil, ErrFullNameRequired
	}

	now := time.Now().UTC()

	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		Phone:        phone,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(
		strings.TrimSpace(email),
	)
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func (u *User) UpdateProfile(
	fullName string,
	phone string,
) error {
	fullName = strings.TrimSpace(fullName)
	phone = strings.TrimSpace(phone)

	if fullName == "" {
		return ErrFullNameRequired
	}

	u.FullName = fullName
	u.Phone = phone
	u.UpdatedAt = time.Now().UTC()

	return nil
}

func (u *User) Activate() error {
	if u.IsActive {
		return ErrUserAlreadyActive
	}

	u.IsActive = true
	u.UpdatedAt = time.Now().UTC()

	return nil
}

func (u *User) Deactivate() error {
	if !u.IsActive {
		return ErrUserAlreadyInactive
	}

	u.IsActive = false
	u.UpdatedAt = time.Now().UTC()

	return nil
}
