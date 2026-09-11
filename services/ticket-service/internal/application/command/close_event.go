package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
)

type CloseEventCommand struct {
	ID uuid.UUID
}

type CloseEventHandler struct {
	repository event.Repository
}

func NewCloseEventHandler(repository event.Repository) *CloseEventHandler {
	return &CloseEventHandler{repository: repository}
}

func (h *CloseEventHandler) Handle(ctx context.Context, cmd CloseEventCommand) (*dto.EventDTO, error) {
	e, err := h.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if err := e.Close(); err != nil {
		return nil, err
	}
	if err := h.repository.Update(ctx, e); err != nil {
		return nil, err
	}
	return dto.FromEvent(e), nil
}
