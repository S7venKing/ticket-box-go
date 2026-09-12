package command

import (
	"context"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/organizer"
)

type CreateOrganizerCommand struct{ Name, Email, Phone, Slug string }
type CreateOrganizerHandler struct{ repository organizer.Repository }

func NewCreateOrganizerHandler(repository organizer.Repository) *CreateOrganizerHandler {
	return &CreateOrganizerHandler{repository: repository}
}
func (h *CreateOrganizerHandler) Handle(ctx context.Context, cmd CreateOrganizerCommand) (*dto.OrganizerDTO, error) {
	o, err := organizer.NewOrganizer(cmd.Name, cmd.Email, cmd.Phone, cmd.Slug)
	if err != nil {
		return nil, err
	}
	if err := h.repository.Create(ctx, o); err != nil {
		return nil, err
	}
	return dto.FromOrganizer(o), nil
}
