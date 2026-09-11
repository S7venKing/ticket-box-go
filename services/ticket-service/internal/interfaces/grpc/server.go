package grpc

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	ticketv1 "github.com/S7venKing/ticket-box-go/gen/ticket/v1"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/command"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/query"
)

type TicketServer struct {
	ticketv1.UnimplementedTicketServiceServer
	logger *slog.Logger

	createOrganizerHandler  *command.CreateOrganizerHandler
	getOrganizerByIDHandler *query.GetOrganizerByIDHandler
	listOrganizersHandler   *query.ListOrganizersHandler
	updateOrganizerHandler  *command.UpdateOrganizerHandler
	createEventHandler      *command.CreateEventHandler
	getEventByIDHandler     *query.GetEventByIDHandler
	listEventsHandler       *query.ListEventsHandler
	publishEventHandler     *command.PublishEventHandler
	submitEventHandler      *command.SubmitEventHandler
	approveEventHandler     *command.ApproveEventHandler
}

func NewTicketServer(
	logger *slog.Logger,
	createOrganizerHandler *command.CreateOrganizerHandler,
	getOrganizerByIDHandler *query.GetOrganizerByIDHandler,
	listOrganizersHandler *query.ListOrganizersHandler,
	updateOrganizerHandler *command.UpdateOrganizerHandler,
	createEventHandler *command.CreateEventHandler,
	getEventByIDHandler *query.GetEventByIDHandler,
	listEventsHandler *query.ListEventsHandler,
	publishEventHandler *command.PublishEventHandler,
	submitEventHandler *command.SubmitEventHandler,
	approveEventHandler *command.ApproveEventHandler,
) *TicketServer {
	return &TicketServer{
		logger:                  logger,
		createOrganizerHandler:  createOrganizerHandler,
		getOrganizerByIDHandler: getOrganizerByIDHandler,
		listOrganizersHandler:   listOrganizersHandler,
		updateOrganizerHandler:  updateOrganizerHandler,
		createEventHandler:      createEventHandler,
		getEventByIDHandler:     getEventByIDHandler,
		listEventsHandler:       listEventsHandler,
		publishEventHandler:     publishEventHandler,
		submitEventHandler:      submitEventHandler,
		approveEventHandler:     approveEventHandler,
	}
}

func (s *TicketServer) CreateOrganizer(ctx context.Context, req *ticketv1.CreateOrganizerRequest) (*ticketv1.CreateOrganizerResponse, error) {
	result, err := s.createOrganizerHandler.Handle(ctx, command.CreateOrganizerCommand{
		Name:  req.GetName(),
		Email: req.GetEmail(),
		Phone: req.GetPhone(),
		Slug:  req.GetSlug(),
	})
	if err != nil {
		return nil, s.mapError(err)
	}
	return &ticketv1.CreateOrganizerResponse{Organizer: toProtoOrganizer(result)}, nil
}

func (s *TicketServer) GetOrganizerById(ctx context.Context, req *ticketv1.GetOrganizerByIdRequest) (*ticketv1.GetOrganizerByIdResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid organizer id")
	}
	result, err := s.getOrganizerByIDHandler.Handle(ctx, query.GetOrganizerByIDQuery{ID: id})
	if err != nil {
		return nil, s.mapError(err)
	}
	return &ticketv1.GetOrganizerByIdResponse{Organizer: toProtoOrganizer(result)}, nil
}

func (s *TicketServer) ListOrganizers(ctx context.Context, req *ticketv1.ListOrganizersRequest) (*ticketv1.ListOrganizersResponse, error) {
	items, err := s.listOrganizersHandler.Handle(ctx, query.ListOrganizersQuery{Offset: int(req.GetOffset()), Limit: int(req.GetLimit())})
	if err != nil {
		return nil, s.mapError(err)
	}
	result := make([]*ticketv1.Organizer, 0, len(items))
	for _, item := range items {
		result = append(result, toProtoOrganizer(item))
	}
	return &ticketv1.ListOrganizersResponse{Organizers: result}, nil
}

func (s *TicketServer) UpdateOrganizer(ctx context.Context, req *ticketv1.UpdateOrganizerRequest) (*ticketv1.UpdateOrganizerResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid organizer id")
	}
	result, err := s.updateOrganizerHandler.Handle(ctx, command.UpdateOrganizerCommand{
		ID:    id,
		Name:  req.GetName(),
		Phone: req.GetPhone(),
		Slug:  req.GetSlug(),
	})
	if err != nil {
		return nil, s.mapError(err)
	}
	return &ticketv1.UpdateOrganizerResponse{Organizer: toProtoOrganizer(result)}, nil
}

func (s *TicketServer) CreateEvent(ctx context.Context, req *ticketv1.CreateEventRequest) (*ticketv1.CreateEventResponse, error) {
	organizerID, err := uuid.Parse(req.GetOrganizerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid organizer id")
	}
	result, err := s.createEventHandler.Handle(ctx, command.CreateEventCommand{
		OrganizerID: organizerID,
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Venue:       req.GetVenue(),
		StartAt:     req.GetStartAt().AsTime(),
		EndAt:       req.GetEndAt().AsTime(),
		Capacity:    int(req.GetCapacity()),
	})
	if err != nil {
		return nil, s.mapError(err)
	}
	return &ticketv1.CreateEventResponse{Event: toProtoEvent(result)}, nil
}

func (s *TicketServer) GetEventById(ctx context.Context, req *ticketv1.GetEventByIdRequest) (*ticketv1.GetEventByIdResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid event id")
	}
	result, err := s.getEventByIDHandler.Handle(ctx, query.GetEventByIDQuery{ID: id})
	if err != nil {
		return nil, s.mapError(err)
	}
	return &ticketv1.GetEventByIdResponse{Event: toProtoEvent(result)}, nil
}

func (s *TicketServer) ListEvents(ctx context.Context, req *ticketv1.ListEventsRequest) (*ticketv1.ListEventsResponse, error) {
	items, err := s.listEventsHandler.Handle(ctx, query.ListEventsQuery{Offset: int(req.GetOffset()), Limit: int(req.GetLimit())})
	if err != nil {
		return nil, s.mapError(err)
	}
	result := make([]*ticketv1.Event, 0, len(items))
	for _, item := range items {
		result = append(result, toProtoEvent(item))
	}
	return &ticketv1.ListEventsResponse{Events: result}, nil
}

func (s *TicketServer) PublishEvent(ctx context.Context, req *ticketv1.PublishEventRequest) (*ticketv1.PublishEventResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid event id")
	}

	result, err := s.publishEventHandler.Handle(ctx, command.PublishEventCommand{ID: id})
	if err != nil {
		return nil, s.mapError(err)
	}

	return &ticketv1.PublishEventResponse{Event: toProtoEvent(result)}, nil
}

func (s *TicketServer) SubmitEvent(ctx context.Context, req *ticketv1.SubmitEventRequest) (*ticketv1.SubmitEventResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid event id")
	}
	result, err := s.submitEventHandler.Handle(ctx, command.SubmitEventCommand{ID: id})
	if err != nil {
		return nil, s.mapError(err)
	}
	return &ticketv1.SubmitEventResponse{Event: toProtoEvent(result)}, nil
}

func (s *TicketServer) ApproveEvent(ctx context.Context, req *ticketv1.ApproveEventRequest) (*ticketv1.ApproveEventResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid event id")
	}
	result, err := s.approveEventHandler.Handle(ctx, command.ApproveEventCommand{ID: id})
	if err != nil {
		return nil, s.mapError(err)
	}
	return &ticketv1.ApproveEventResponse{Event: toProtoEvent(result)}, nil
}

func (s *TicketServer) mapError(err error) error {
	if err == nil {
		return nil
	}
	switch err {
	case status.Error(codes.InvalidArgument, ""):
		return status.Error(codes.InvalidArgument, err.Error())
	}
	// domain-level errors are normalized here so service API stays consistent.
	return status.Error(codes.InvalidArgument, err.Error())
}

func parseProtoTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
