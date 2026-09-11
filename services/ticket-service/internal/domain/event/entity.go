package event

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusDraft           Status = "draft"
	StatusPendingApproval Status = "pending_approval"
	StatusPublished       Status = "published"
	StatusCancelled       Status = "cancelled"
	StatusClosed          Status = "closed"
)

type Event struct {
	ID          uuid.UUID
	OrganizerID uuid.UUID
	Title       string
	Description string
	Venue       string
	StartAt     time.Time
	EndAt       time.Time
	Capacity    int
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewEvent(organizerID uuid.UUID, title, description, venue string, startAt, endAt time.Time, capacity int) (*Event, error) {
	if organizerID == uuid.Nil {
		return nil, ErrOrganizerIDRequired
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrTitleRequired
	}
	venue = strings.TrimSpace(venue)
	if venue == "" {
		return nil, ErrVenueRequired
	}
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}
	if !endAt.After(startAt) {
		return nil, ErrInvalidDateRange
	}

	now := time.Now().UTC()
	return &Event{
		ID:          uuid.New(),
		OrganizerID: organizerID,
		Title:       title,
		Description: description,
		Venue:       venue,
		StartAt:     startAt.UTC(),
		EndAt:       endAt.UTC(),
		Capacity:    capacity,
		Status:      StatusDraft,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (e *Event) Publish() error {
	if e.Status == StatusPublished {
		return ErrPublishedEvent
	}
	if e.Status != StatusDraft && e.Status != StatusPendingApproval {
		return ErrNotPendingApprovalEvent
	}

	e.Status = StatusPublished
	e.UpdatedAt = time.Now().UTC()
	return nil
}

func (e *Event) SubmitForApproval() error {
	if e.Status != StatusDraft {
		return ErrNotDraftEvent
	}
	e.Status = StatusPendingApproval
	e.UpdatedAt = time.Now().UTC()
	return nil
}

func (e *Event) Cancel() error {
	if e.Status == StatusCancelled {
		return ErrCancelledEvent
	}
	e.Status = StatusCancelled
	e.UpdatedAt = time.Now().UTC()
	return nil
}

func (e *Event) Close() error {
	if e.Status == StatusClosed {
		return ErrInvalidStatus
	}
	e.Status = StatusClosed
	e.UpdatedAt = time.Now().UTC()
	return nil
}

func (e *Event) UpdateDetails(title, description, venue string, startAt, endAt time.Time, capacity int) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return ErrTitleRequired
	}
	venue = strings.TrimSpace(venue)
	if venue == "" {
		return ErrVenueRequired
	}
	if capacity <= 0 {
		return ErrInvalidCapacity
	}
	if !endAt.After(startAt) {
		return ErrInvalidDateRange
	}

	e.Title = title
	e.Description = description
	e.Venue = venue
	e.StartAt = startAt.UTC()
	e.EndAt = endAt.UTC()
	e.Capacity = capacity
	e.UpdatedAt = time.Now().UTC()
	return nil
}
