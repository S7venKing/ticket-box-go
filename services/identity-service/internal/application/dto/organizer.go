package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/organizer"
)

type OrganizerDTO struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Phone     string
	Slug      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FromOrganizer(o *organizer.Organizer) *OrganizerDTO {
	if o == nil {
		return nil
	}
	return &OrganizerDTO{ID: o.ID, Name: o.Name, Email: o.Email, Phone: o.Phone, Slug: o.Slug, IsActive: o.IsActive, CreatedAt: o.CreatedAt, UpdatedAt: o.UpdatedAt}
}
