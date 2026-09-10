package command

import (
	"context"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
)

type AuthenticateUserCommand struct {
	Email    string
	Password string
}

type AuthenticateUserHandler struct {
	userRepository user.Repository
	passwordHasher PasswordHasher
}

func NewAuthenticateUserHandler(
	userRepository user.Repository,
	passwordHasher PasswordHasher,
) *AuthenticateUserHandler {
	return &AuthenticateUserHandler{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
	}
}

func (h *AuthenticateUserHandler) Handle(
	ctx context.Context,
	cmd AuthenticateUserCommand,
) (*dto.UserDTO, error) {
	if cmd.Password == "" {
		return nil, ErrPasswordRequired
	}

	email := user.NormalizeEmail(cmd.Email)
	if err := user.ValidateEmail(email); err != nil {
		return nil, err
	}

	u, err := h.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !u.IsActive {
		return nil, user.ErrUserInactive
	}

	if err := h.passwordHasher.Compare(cmd.Password, u.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	return dto.FromUser(u), nil
}
