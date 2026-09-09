package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
)

type UpdateProfileCommand struct {
	UserID   uuid.UUID
	FullName string
	Phone    string
}

type UpdateProfileHandler struct {
	userRepository user.Repository
}

func NewUpdateProfileHandler(
	userRepository user.Repository,
) *UpdateProfileHandler {
	return &UpdateProfileHandler{
		userRepository: userRepository,
	}
}

func (h *UpdateProfileHandler) Handle(
	ctx context.Context,
	cmd UpdateProfileCommand,
) (*dto.UserDTO, error) {
	u, err := h.userRepository.GetByID(
		ctx,
		cmd.UserID,
	)
	if err != nil {
		return nil, err
	}

	if err := u.UpdateProfile(
		cmd.FullName,
		cmd.Phone,
	); err != nil {
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
