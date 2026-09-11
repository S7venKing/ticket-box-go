package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/organizer"
)

type ActivateOrganizerCommand struct {
	ID uuid.UUID
}

type ActivateOrganizerHandler struct {
	repository organizer.Repository
}

func NewActivateOrganizerHandler(repository organizer.Repository) *ActivateOrganizerHandler {
	return &ActivateOrganizerHandler{repository: repository}
}

func (h *ActivateOrganizerHandler) Handle(ctx context.Context, cmd ActivateOrganizerCommand) (*dto.OrganizerDTO, error) {
	o, err := h.repository.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if err := o.Activate(); err != nil {
		return nil, err
	}
	if err := h.repository.Update(ctx, o); err != nil {
		return nil, err
	}
	return dto.FromOrganizer(o), nil
}
