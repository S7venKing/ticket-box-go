package query

import (
	"context"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/organizer"
)

type GetOrganizerByEmailQuery struct {
	Email string
}

type GetOrganizerByEmailHandler struct {
	repository organizer.Repository
}

func NewGetOrganizerByEmailHandler(repository organizer.Repository) *GetOrganizerByEmailHandler {
	return &GetOrganizerByEmailHandler{repository: repository}
}

func (h *GetOrganizerByEmailHandler) Handle(ctx context.Context, query GetOrganizerByEmailQuery) (*dto.OrganizerDTO, error) {
	o, err := h.repository.GetByEmail(ctx, query.Email)
	if err != nil {
		return nil, err
	}
	return dto.FromOrganizer(o), nil
}
