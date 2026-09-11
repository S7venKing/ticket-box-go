package organizer

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, organizer *Organizer) error
	GetByID(ctx context.Context, id uuid.UUID) (*Organizer, error)
	GetByEmail(ctx context.Context, email string) (*Organizer, error)
	Update(ctx context.Context, organizer *Organizer) error
	List(ctx context.Context, offset, limit int) ([]*Organizer, error)
}
