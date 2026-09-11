package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/organizer"
)

type UpdateOrganizerCommand struct {
	ID    uuid.UUID
	Name  string
	Phone string
	Slug  string
}

type UpdateOrganizerHandler struct {
	repository organizer.Repository
}

func NewUpdateOrganizerHandler(repository organizer.Repository) *UpdateOrganizerHandler {
	return &UpdateOrganizerHandler{repository: repository}
}

func (h *UpdateOrganizerHandler) Handle(ctx context.Context, cmd UpdateOrganizerCommand) (*dto.OrganizerDTO, error) {
	o, err := h.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if err := o.UpdateProfile(cmd.Name, cmd.Phone, cmd.Slug); err != nil {
		return nil, err
	}

	if err := h.repository.Update(ctx, o); err != nil {
		return nil, err
	}

	return dto.FromOrganizer(o), nil
}
