package organizer

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(context.Context, *Organizer) error
	GetByID(context.Context, uuid.UUID) (*Organizer, error)
	Update(context.Context, *Organizer) error
	List(context.Context, int, int) ([]*Organizer, error)
}
