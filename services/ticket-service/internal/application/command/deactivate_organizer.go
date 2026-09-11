package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/organizer"
)

type DeactivateOrganizerCommand struct {
	ID uuid.UUID
}

type DeactivateOrganizerHandler struct {
	repository organizer.Repository
}

func NewDeactivateOrganizerHandler(repository organizer.Repository) *DeactivateOrganizerHandler {
	return &DeactivateOrganizerHandler{repository: repository}
}

func (h *DeactivateOrganizerHandler) Handle(ctx context.Context, cmd DeactivateOrganizerCommand) (*dto.OrganizerDTO, error) {
	o, err := h.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if err := o.Deactivate(); err != nil {
		return nil, err
	}
	if err := h.repository.Update(ctx, o); err != nil {
		return nil, err
	}
	return dto.FromOrganizer(o), nil
}
