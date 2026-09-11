package event_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
)

func TestNewEvent_Valid(t *testing.T) {
	start := time.Now().Add(24 * time.Hour)
	end := start.Add(2 * time.Hour)

	e, err := event.NewEvent(uuid.New(), "Go Conference", "Annual meetup", "HCM", start, end, 200)
	require.NoError(t, err)
	require.Equal(t, event.StatusDraft, e.Status)
	require.Equal(t, 200, e.Capacity)
	require.True(t, e.EndAt.After(e.StartAt))
}

func TestEvent_Publish(t *testing.T) {
	start := time.Now().Add(24 * time.Hour)
	end := start.Add(2 * time.Hour)
	e, err := event.NewEvent(uuid.New(), "Launch Event", "Desc", "HCM", start, end, 100)
	require.NoError(t, err)
	require.NoError(t, e.Publish())
	require.Equal(t, event.StatusPublished, e.Status)
	require.ErrorIs(t, e.Publish(), event.ErrPublishedEvent)
}

func TestEvent_Cancel(t *testing.T) {
	start := time.Now().Add(24 * time.Hour)
	end := start.Add(2 * time.Hour)
	e, err := event.NewEvent(uuid.New(), "Launch Event", "Desc", "HCM", start, end, 100)
	require.NoError(t, err)
	require.NoError(t, e.Cancel())
	require.Equal(t, event.StatusCancelled, e.Status)
	require.ErrorIs(t, e.Cancel(), event.ErrCancelledEvent)
}

func TestEvent_InvalidDateRange(t *testing.T) {
	start := time.Now().Add(2 * time.Hour)
	end := start.Add(-1 * time.Hour)
	_, err := event.NewEvent(uuid.New(), "Bad Event", "Desc", "HCM", start, end, 10)
	require.ErrorIs(t, err, event.ErrInvalidDateRange)
}
