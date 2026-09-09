package user

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User is the identity aggregate root.
// PasswordHash never leaves the application boundary: dto.UserDTO and the
// protobuf User message do not carry it.
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
	email = NormalizeEmail(email)
	fullName = strings.TrimSpace(fullName)
	phone = strings.TrimSpace(phone)

	if err := ValidateEmail(email); err != nil {
		return nil, err
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

// NormalizeEmail returns the canonical form used for storage and lookups.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// ValidateEmail checks a normalized email address.
func ValidateEmail(email string) error {
	if email == "" {
		return ErrEmailRequired
	}

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		// addr.Address != email rejects display-name forms like "Name <a@b.c>".
		return ErrInvalidEmail
	}

	return nil
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
