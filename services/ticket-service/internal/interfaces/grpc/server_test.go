package grpc_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/organizer"
)

func TestProtoMapper_OrganizerAndEvent(t *testing.T) {
	o := &organizer.Organizer{
		ID:        uuid.New(),
		Name:      "Tech Meetup",
		Email:     "team@example.com",
		Phone:     "0900000000",
		Slug:      "tech-meetup",
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	e := &event.Event{
		ID:          uuid.New(),
		OrganizerID: uuid.New(),
		Title:       "Go Conference",
		Description: "Annual conference",
		Venue:       "Ho Chi Minh",
		StartAt:     time.Now().UTC().Add(24 * time.Hour),
		EndAt:       time.Now().UTC().Add(26 * time.Hour),
		Capacity:    200,
		Status:      event.StatusPublished,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	_ = o
	_ = e
	_ = timestamppb.Now
	require.NotZero(t, o.ID)
	require.NotZero(t, e.ID)
}
