package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
)

type PublishEventCommand struct {
	ID uuid.UUID
}

type PublishEventHandler struct {
	repository event.Repository
}

func NewPublishEventHandler(repository event.Repository) *PublishEventHandler {
	return &PublishEventHandler{repository: repository}
}

func (h *PublishEventHandler) Handle(ctx context.Context, cmd PublishEventCommand) (*dto.EventDTO, error) {
	e, err := h.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if err := e.Publish(); err != nil {
		return nil, err
	}
	if err := h.repository.Update(ctx, e); err != nil {
		return nil, err
	}
	return dto.FromEvent(e), nil
}
