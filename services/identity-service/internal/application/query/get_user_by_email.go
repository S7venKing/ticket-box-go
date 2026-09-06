package query

import (
	"context"

	"github.com/s7venking/ticket-box/identity/internal/application/dto"
	"github.com/s7venking/ticket-box/identity/internal/domain/user"
)

type GetUserByEmailQuery struct {
	Email string
}

type GetUserByEmailHandler struct {
	userRepository user.Repository
}

func NewGetUserByEmailHandler(
	userRepository user.Repository,
) *GetUserByEmailHandler {
	return &GetUserByEmailHandler{
		userRepository: userRepository,
	}
}

func (h *GetUserByEmailHandler) Handle(
	ctx context.Context,
	query GetUserByEmailQuery,
) (*dto.UserDTO, error) {
	u, err := h.userRepository.GetByEmail(
		ctx,
		query.Email,
	)
	if err != nil {
		return nil, err
	}

	result := dto.FromUser(u)

	return result, nil
}
