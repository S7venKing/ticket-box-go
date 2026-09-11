package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
)

type CancelEventCommand struct {
	ID uuid.UUID
}

type CancelEventHandler struct {
	repository event.Repository
}

func NewCancelEventHandler(repository event.Repository) *CancelEventHandler {
	return &CancelEventHandler{repository: repository}
}

func (h *CancelEventHandler) Handle(ctx context.Context, cmd CancelEventCommand) (*dto.EventDTO, error) {
	e, err := h.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if err := e.Cancel(); err != nil {
		return nil, err
	}
	if err := h.repository.Update(ctx, e); err != nil {
		return nil, err
	}
	return dto.FromEvent(e), nil
}
