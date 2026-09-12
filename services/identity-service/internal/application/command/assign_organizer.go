package command

import (
	"context"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
	"github.com/google/uuid"
	"time"
)

type AssignUserToOrganizerCommand struct{ UserID, OrganizerID uuid.UUID }

type AssignUserToOrganizerHandler struct{ repo user.Repository }

func NewAssignUserToOrganizerHandler(repo user.Repository) *AssignUserToOrganizerHandler {
	return &AssignUserToOrganizerHandler{repo: repo}
}
func (h *AssignUserToOrganizerHandler) Handle(ctx context.Context, cmd AssignUserToOrganizerCommand) (*dto.UserDTO, error) {
	u, err := h.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	u.OrganizerID = &cmd.OrganizerID
	u.UpdatedAt = time.Now().UTC()
	if err := h.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	return dto.FromUser(u), nil
}
