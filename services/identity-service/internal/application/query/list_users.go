package query

import (
	"context"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
)

type ListUsersHandler struct {
	repository user.ListRepository
}

func NewListUsersHandler(repository user.ListRepository) *ListUsersHandler {
	return &ListUsersHandler{repository: repository}
}

func (h *ListUsersHandler) Handle(ctx context.Context, offset, limit int) ([]*user.User, error) {
	return h.repository.List(ctx, offset, limit)
}
