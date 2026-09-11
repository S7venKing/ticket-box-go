package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
)

type EventDTO struct {
	ID          uuid.UUID
	OrganizerID uuid.UUID
	Title       string
	Description string
	Venue       string
	StartAt     time.Time
	EndAt       time.Time
	Capacity    int
	Status      event.Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func FromEvent(e *event.Event) *EventDTO {
	if e == nil {
		return nil
	}

	return &EventDTO{
		ID:          e.ID,
		OrganizerID: e.OrganizerID,
		Title:       e.Title,
		Description: e.Description,
		Venue:       e.Venue,
		StartAt:     e.StartAt,
		EndAt:       e.EndAt,
		Capacity:    e.Capacity,
		Status:      e.Status,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
