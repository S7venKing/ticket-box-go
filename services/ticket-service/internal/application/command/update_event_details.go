package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
)

type UpdateEventDetailsCommand struct {
	ID          uuid.UUID
	Title       string
	Description string
	Venue       string
	StartAt     time.Time
	EndAt       time.Time
	Capacity    int
}

type UpdateEventDetailsHandler struct {
	repository event.Repository
}

func NewUpdateEventDetailsHandler(repository event.Repository) *UpdateEventDetailsHandler {
	return &UpdateEventDetailsHandler{repository: repository}
}

func (h *UpdateEventDetailsHandler) Handle(ctx context.Context, cmd UpdateEventDetailsCommand) (*dto.EventDTO, error) {
	e, err := h.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if err := e.UpdateDetails(cmd.Title, cmd.Description, cmd.Venue, cmd.StartAt, cmd.EndAt, cmd.Capacity); err != nil {
		return nil, err
	}
	if err := h.repository.Update(ctx, e); err != nil {
		return nil, err
	}
	return dto.FromEvent(e), nil
}
