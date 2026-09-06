package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/s7venking/ticket-box/identity/internal/domain/user"
)

type UserDTO struct {
	ID        uuid.UUID
	Email     string
	FullName  string
	Phone     string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FromUser(u *user.User) *UserDTO {
	return &UserDTO{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		Phone:     u.Phone,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
