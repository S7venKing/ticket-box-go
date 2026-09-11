package command

import (
	"context"
	"unicode/utf8"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
)

type CreateUserCommand struct {
	Email    string
	Password string
	FullName string
	Phone    string
	Role     string
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
	// Cheap validations first: no DB round-trip and no argon2 work for bad input.
	if err := validatePassword(cmd.Password); err != nil {
		return nil, err
	}

	email := user.NormalizeEmail(cmd.Email)

	if err := user.ValidateEmail(email); err != nil {
		return nil, err
	}

	exists, err := h.userRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, user.ErrEmailAlreadyExists
	}

	passwordHash, err := h.passwordHasher.Hash(cmd.Password)
	if err != nil {
		return nil, err
	}

	u, err := user.NewUser(email, passwordHash, cmd.FullName, cmd.Phone, cmd.Role)
	if err != nil {
		return nil, err
	}

	// The UNIQUE index on users.email is the real guard against concurrent
	// creates; the repository maps a duplicate-key violation to
	// user.ErrEmailAlreadyExists.
	if err := h.userRepository.Create(ctx, u); err != nil {
		return nil, err
	}

	return dto.FromUser(u), nil
}

func validatePassword(password string) error {
	if password == "" {
		return ErrPasswordRequired
	}

	if utf8.RuneCountInString(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}

	return nil
}
