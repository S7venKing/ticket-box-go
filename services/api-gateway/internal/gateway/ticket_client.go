package gateway

import (
	"context"
	"time"

	ticketv1 "github.com/S7venKing/ticket-box-go/gen/ticket/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TicketClient struct {
	client ticketv1.TicketServiceClient
}

func NewTicketClient(cc grpc.ClientConnInterface) *TicketClient {
	return &TicketClient{client: ticketv1.NewTicketServiceClient(cc)}
}

func (c *TicketClient) CreateEvent(ctx context.Context, organizerID, title, description, venue string, startAt, endAt time.Time, capacity int32) (*ticketv1.Event, error) {
	resp, err := c.client.CreateEvent(ctx, &ticketv1.CreateEventRequest{
		OrganizerId: organizerID,
		Title:       title,
		Description: description,
		Venue:       venue,
		StartAt:     timestamptz(startAt),
		EndAt:       timestamptz(endAt),
		Capacity:    capacity,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetEvent(), nil
}

func (c *TicketClient) GetEventByID(ctx context.Context, id string) (*ticketv1.Event, error) {
	resp, err := c.client.GetEventById(ctx, &ticketv1.GetEventByIdRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return resp.GetEvent(), nil
}

func (c *TicketClient) ListEvents(ctx context.Context, offset, limit int32) ([]*ticketv1.Event, error) {
	resp, err := c.client.ListEvents(ctx, &ticketv1.ListEventsRequest{Offset: offset, Limit: limit})
	if err != nil {
		return nil, err
	}
	return resp.GetEvents(), nil
}

func (c *TicketClient) PublishEvent(ctx context.Context, id string) (*ticketv1.Event, error) {
	resp, err := c.client.PublishEvent(ctx, &ticketv1.PublishEventRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return resp.GetEvent(), nil
}

func (c *TicketClient) SubmitEvent(ctx context.Context, id string) (*ticketv1.Event, error) {
	resp, err := c.client.SubmitEvent(ctx, &ticketv1.SubmitEventRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return resp.GetEvent(), nil
}

func (c *TicketClient) ApproveEvent(ctx context.Context, id string) (*ticketv1.Event, error) {
	resp, err := c.client.ApproveEvent(ctx, &ticketv1.ApproveEventRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return resp.GetEvent(), nil
}

func timestamptz(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t.UTC())
}
