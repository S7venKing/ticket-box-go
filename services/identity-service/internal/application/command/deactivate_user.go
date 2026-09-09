package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
)

type DeactivateUserCommand struct {
	UserID uuid.UUID
}

type DeactivateUserHandler struct {
	userRepository user.Repository
}

func NewDeactivateUserHandler(
	userRepository user.Repository,
) *DeactivateUserHandler {
	return &DeactivateUserHandler{
		userRepository: userRepository,
	}
}

func (h *DeactivateUserHandler) Handle(
	ctx context.Context,
	cmd DeactivateUserCommand,
) (*dto.UserDTO, error) {
	u, err := h.userRepository.GetByID(
		ctx,
		cmd.UserID,
	)
	if err != nil {
		return nil, err
	}

	if err := u.Deactivate(); err != nil {
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
