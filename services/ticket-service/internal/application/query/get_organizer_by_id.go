package query

import (
	"context"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/organizer"
)

type GetOrganizerByIDQuery struct {
	ID uuid.UUID
}

type GetOrganizerByIDHandler struct {
	repository organizer.Repository
}

func NewGetOrganizerByIDHandler(repository organizer.Repository) *GetOrganizerByIDHandler {
	return &GetOrganizerByIDHandler{repository: repository}
}

func (h *GetOrganizerByIDHandler) Handle(ctx context.Context, query GetOrganizerByIDQuery) (*dto.OrganizerDTO, error) {
	o, err := h.repository.GetByID(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	return dto.FromOrganizer(o), nil
}
