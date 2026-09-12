package query

import (
	"context"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/organizer"
	"github.com/google/uuid"
)

type GetOrganizerHandler struct{ repository organizer.Repository }

func NewGetOrganizerHandler(repository organizer.Repository) *GetOrganizerHandler {
	return &GetOrganizerHandler{repository: repository}
}
func (h *GetOrganizerHandler) Handle(ctx context.Context, id uuid.UUID) (*dto.OrganizerDTO, error) {
	o, err := h.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.FromOrganizer(o), nil
}
