package event

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, event *Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*Event, error)
	GetByOrganizer(ctx context.Context, organizerID uuid.UUID, offset, limit int) ([]*Event, error)
	Update(ctx context.Context, event *Event) error
	List(ctx context.Context, offset, limit int) ([]*Event, error)
}
