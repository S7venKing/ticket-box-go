package query

import (
	"context"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
)

type ListEventsQuery struct {
	Offset int
	Limit  int
}

type ListEventsHandler struct {
	repository event.Repository
}

func NewListEventsHandler(repository event.Repository) *ListEventsHandler {
	return &ListEventsHandler{repository: repository}
}

func (h *ListEventsHandler) Handle(ctx context.Context, query ListEventsQuery) ([]*dto.EventDTO, error) {
	items, err := h.repository.List(ctx, query.Offset, query.Limit)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.EventDTO, 0, len(items))
	for _, item := range items {
		result = append(result, dto.FromEvent(item))
	}
	return result, nil
}
