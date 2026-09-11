package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
)

type GetEventByIDQuery struct {
	ID uuid.UUID
}

type GetEventByIDHandler struct {
	repository event.Repository
}

func NewGetEventByIDHandler(repository event.Repository) *GetEventByIDHandler {
	return &GetEventByIDHandler{repository: repository}
}

func (h *GetEventByIDHandler) Handle(ctx context.Context, query GetEventByIDQuery) (*dto.EventDTO, error) {
	e, err := h.repository.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	return dto.FromEvent(e), nil
}
