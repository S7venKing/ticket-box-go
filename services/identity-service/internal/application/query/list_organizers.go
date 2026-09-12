package query

import (
	"context"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/organizer"
)

type ListOrganizersHandler struct{ repository organizer.Repository }

func NewListOrganizersHandler(repository organizer.Repository) *ListOrganizersHandler {
	return &ListOrganizersHandler{repository: repository}
}
func (h *ListOrganizersHandler) Handle(ctx context.Context, offset, limit int) ([]*dto.OrganizerDTO, error) {
	items, err := h.repository.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.OrganizerDTO, 0, len(items))
	for _, item := range items {
		result = append(result, dto.FromOrganizer(item))
	}
	return result, nil
}
