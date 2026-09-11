package event

import "errors"

var (
	ErrOrganizerIDRequired     = errors.New("organizer id is required")
	ErrTitleRequired           = errors.New("event title is required")
	ErrVenueRequired           = errors.New("event venue is required")
	ErrCapacityRequired        = errors.New("event capacity is required")
	ErrInvalidCapacity         = errors.New("event capacity must be greater than zero")
	ErrInvalidDateRange        = errors.New("event end time must be after start time")
	ErrInvalidStatus           = errors.New("event status is invalid")
	ErrPublishedEvent          = errors.New("event is already published")
	ErrCancelledEvent          = errors.New("event is already cancelled")
	ErrNotDraftEvent           = errors.New("only draft events can be published")
	ErrNotPendingApprovalEvent = errors.New("only pending approval events can be published")
	ErrEventNotFound           = errors.New("event not found")
)
