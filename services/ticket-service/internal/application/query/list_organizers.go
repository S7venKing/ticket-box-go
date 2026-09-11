package query

import (
	"context"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/organizer"
)

type ListOrganizersQuery struct {
	Offset int
	Limit  int
}

type ListOrganizersHandler struct {
	repository organizer.Repository
}

func NewListOrganizersHandler(repository organizer.Repository) *ListOrganizersHandler {
	return &ListOrganizersHandler{repository: repository}
}

func (h *ListOrganizersHandler) Handle(ctx context.Context, query ListOrganizersQuery) ([]*dto.OrganizerDTO, error) {
	items, err := h.repository.List(ctx, query.Offset, query.Limit)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.OrganizerDTO, 0, len(items))
	for _, item := range items {
		result = append(result, dto.FromOrganizer(item))
	}
	return result, nil
}
