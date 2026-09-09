package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
)

type GetUserByIDQuery struct {
	ID uuid.UUID
}

type GetUserByIDHandler struct {
	userRepository user.Repository
}

func NewGetUserByIDHandler(
	userRepository user.Repository,
) *GetUserByIDHandler {
	return &GetUserByIDHandler{
		userRepository: userRepository,
	}
}

func (h *GetUserByIDHandler) Handle(
	ctx context.Context,
	query GetUserByIDQuery,
) (*dto.UserDTO, error) {
	u, err := h.userRepository.GetByID(
		ctx,
		query.ID,
	)
	if err != nil {
		return nil, err
	}

	result := dto.FromUser(u)

	return result, nil
}
