package organizer

import "errors"

var (
	ErrNameRequired             = errors.New("organizer name is required")
	ErrEmailRequired            = errors.New("organizer email is required")
	ErrInvalidEmail             = errors.New("organizer email is invalid")
	ErrSlugRequired             = errors.New("organizer slug is required")
	ErrInvalidSlug              = errors.New("organizer slug is invalid")
	ErrOrganizerInactive        = errors.New("organizer is inactive")
	ErrOrganizerAlreadyActive   = errors.New("organizer is already active")
	ErrOrganizerAlreadyInactive = errors.New("organizer is already inactive")
	ErrOrganizerNotFound        = errors.New("organizer not found")
	ErrEmailAlreadyExists       = errors.New("organizer email already exists")
)
