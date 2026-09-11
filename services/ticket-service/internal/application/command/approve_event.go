package command

import (
	"context"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
	"github.com/google/uuid"
)

type ApproveEventCommand struct{ ID uuid.UUID }

type ApproveEventHandler struct{ repository event.Repository }

func NewApproveEventHandler(repository event.Repository) *ApproveEventHandler {
	return &ApproveEventHandler{repository: repository}
}

func (h *ApproveEventHandler) Handle(ctx context.Context, cmd ApproveEventCommand) (*dto.EventDTO, error) {
	item, err := h.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if err := item.Publish(); err != nil {
		return nil, err
	}
	if err := h.repository.Update(ctx, item); err != nil {
		return nil, err
	}
	return dto.FromEvent(item), nil
}
