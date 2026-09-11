package grpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	ticketv1 "github.com/S7venKing/ticket-box-go/gen/ticket/v1"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/dto"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
)

func toProtoOrganizer(o *dto.OrganizerDTO) *ticketv1.Organizer {
	if o == nil {
		return nil
	}
	return &ticketv1.Organizer{
		Id:        o.ID.String(),
		Name:      o.Name,
		Email:     o.Email,
		Phone:     o.Phone,
		Slug:      o.Slug,
		IsActive:  o.IsActive,
		CreatedAt: timestamppb.New(o.CreatedAt),
		UpdatedAt: timestamppb.New(o.UpdatedAt),
	}
}

func toProtoEvent(e *dto.EventDTO) *ticketv1.Event {
	if e == nil {
		return nil
	}
	return &ticketv1.Event{
		Id:          e.ID.String(),
		OrganizerId: e.OrganizerID.String(),
		Title:       e.Title,
		Description: e.Description,
		Venue:       e.Venue,
		StartAt:     timestamppb.New(e.StartAt),
		EndAt:       timestamppb.New(e.EndAt),
		Capacity:    int32(e.Capacity),
		Status:      toProtoEventStatus(e.Status),
		CreatedAt:   timestamppb.New(e.CreatedAt),
		UpdatedAt:   timestamppb.New(e.UpdatedAt),
	}
}

func toProtoEventStatus(status event.Status) ticketv1.EventStatus {
	switch status {
	case event.StatusDraft:
		return ticketv1.EventStatus_EVENT_STATUS_DRAFT
	case event.StatusPendingApproval:
		return ticketv1.EventStatus_EVENT_STATUS_PENDING_APPROVAL
	case event.StatusPublished:
		return ticketv1.EventStatus_EVENT_STATUS_PUBLISHED
	case event.StatusCancelled:
		return ticketv1.EventStatus_EVENT_STATUS_CANCELLED
	case event.StatusClosed:
		return ticketv1.EventStatus_EVENT_STATUS_CLOSED
	default:
		return ticketv1.EventStatus_EVENT_STATUS_UNSPECIFIED
	}
}
