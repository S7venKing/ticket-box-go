package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
)

type CreateEventCommand struct {
	OrganizerID uuid.UUID
	Title       string
	Description string
	Venue       string
	StartAt     time.Time
	EndAt       time.Time
	Capacity    int
}

type CreateEventHandler struct {
	repository event.Repository
}

func NewCreateEventHandler(repository event.Repository) *CreateEventHandler {
	return &CreateEventHandler{repository: repository}
}

func (h *CreateEventHandler) Handle(ctx context.Context, cmd CreateEventCommand) (*dto.EventDTO, error) {
	e, err := event.NewEvent(cmd.OrganizerID, cmd.Title, cmd.Description, cmd.Venue, cmd.StartAt, cmd.EndAt, cmd.Capacity)
	if err != nil {
		return nil, err
	}

	if err := h.repository.Create(ctx, e); err != nil {
		return nil, err
	}

	return dto.FromEvent(e), nil
}
