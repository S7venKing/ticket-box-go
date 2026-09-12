package organizer

import "errors"

var (
	ErrNameRequired       = errors.New("organizer name is required")
	ErrEmailRequired      = errors.New("organizer email is required")
	ErrInvalidEmail       = errors.New("organizer email is invalid")
	ErrSlugRequired       = errors.New("organizer slug is required")
	ErrInvalidSlug        = errors.New("organizer slug is invalid")
	ErrOrganizerNotFound  = errors.New("organizer not found")
	ErrOrganizerDuplicate = errors.New("organizer email or slug already exists")
)
