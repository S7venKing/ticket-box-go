package command

import (
	"context"
	"errors"

	"github.com/s7venking/ticket-box/identity/internal/application/dto"
	"github.com/s7venking/ticket-box/identity/internal/domain/user"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

type CreateUserCommand struct {
	Email    string
	Password string
	FullName string
	Phone    string
}

type CreateUserHandler struct {
	userRepository user.Repository
	passwordHasher PasswordHasher
}

func NewCreateUserHandler(
	userRepository user.Repository,
	passwordHasher PasswordHasher,
) *CreateUserHandler {
	return &CreateUserHandler{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
	}
}

func (h *CreateUserHandler) Handle(
	ctx context.Context,
	cmd CreateUserCommand,
) (*dto.UserDTO, error) {
	exists, err := h.userRepository.ExistsByEmail(
		ctx,
		cmd.Email,
	)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrEmailAlreadyExists
	}

	passwordHash, err := h.passwordHasher.Hash(
		cmd.Password,
	)
	if err != nil {
		return nil, err
	}

	u, err := user.NewUser(
		cmd.Email,
		passwordHash,
		cmd.FullName,
		cmd.Phone,
	)
	if err != nil {
		return nil, err
	}

	if err := h.userRepository.Create(ctx, u); err != nil {
		return nil, err
	}

	return dto.FromUser(u), nil
}
