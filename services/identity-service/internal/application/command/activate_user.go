package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
)

type ActivateUserCommand struct {
	UserID uuid.UUID
}

type ActivateUserHandler struct {
	userRepository user.Repository
}

func NewActivateUserHandler(
	userRepository user.Repository,
) *ActivateUserHandler {
	return &ActivateUserHandler{
		userRepository: userRepository,
	}
}

func (h *ActivateUserHandler) Handle(
	ctx context.Context,
	cmd ActivateUserCommand,
) (*dto.UserDTO, error) {
	u, err := h.userRepository.GetByID(
		ctx,
		cmd.UserID,
	)
	if err != nil {
		return nil, err
	}

	if err := u.Activate(); err != nil {
		return nil, err
	}

	if err := h.userRepository.Update(
		ctx,
		u,
	); err != nil {
		return nil, err
	}

	result := dto.FromUser(u)

	return result, nil
}
