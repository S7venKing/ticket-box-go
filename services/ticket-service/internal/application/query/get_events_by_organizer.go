package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
)

type GetEventsByOrganizerQuery struct {
	OrganizerID uuid.UUID
	Offset      int
	Limit       int
}

type GetEventsByOrganizerHandler struct {
	repository event.Repository
}

func NewGetEventsByOrganizerHandler(repository event.Repository) *GetEventsByOrganizerHandler {
	return &GetEventsByOrganizerHandler{repository: repository}
}

func (h *GetEventsByOrganizerHandler) Handle(ctx context.Context, query GetEventsByOrganizerQuery) ([]*dto.EventDTO, error) {
	items, err := h.repository.GetByOrganizer(ctx, query.OrganizerID, query.Offset, query.Limit)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.EventDTO, 0, len(items))
	for _, item := range items {
		result = append(result, dto.FromEvent(item))
	}
	return result, nil
}
