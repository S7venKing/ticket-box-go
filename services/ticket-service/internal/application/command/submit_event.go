package command

import (
	"context"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
	"github.com/google/uuid"
)

type SubmitEventCommand struct{ ID uuid.UUID }

type SubmitEventHandler struct{ repository event.Repository }

func NewSubmitEventHandler(repository event.Repository) *SubmitEventHandler {
	return &SubmitEventHandler{repository: repository}
}

func (h *SubmitEventHandler) Handle(ctx context.Context, cmd SubmitEventCommand) (*dto.EventDTO, error) {
	item, err := h.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if err := item.SubmitForApproval(); err != nil {
		return nil, err
	}
	if err := h.repository.Update(ctx, item); err != nil {
		return nil, err
	}
	return dto.FromEvent(item), nil
}
