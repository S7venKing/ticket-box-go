package organizer

import (
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Organizer struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Phone     string
	Slug      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOrganizer(name, email, phone, slug string) (*Organizer, error) {
	name = strings.TrimSpace(name)
	email = NormalizeEmail(email)
	phone = strings.TrimSpace(phone)
	if slug != "" {
		slug = strings.TrimSpace(slug)
	}

	if name == "" {
		return nil, ErrNameRequired
	}
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}
	if slug == "" {
		slug = toSlug(name)
	}
	if err := ValidateSlug(slug); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Organizer{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Slug:      slug,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ValidateEmail(email string) error {
	if email == "" {
		return ErrEmailRequired
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return ErrInvalidEmail
	}
	return nil
}

func ValidateSlug(slug string) error {
	if slug == "" {
		return ErrSlugRequired
	}
	if !slugPattern.MatchString(slug) {
		return ErrInvalidSlug
	}
	return nil
}

func toSlug(value string) string {
	lower := strings.ToLower(strings.TrimSpace(value))
	lower = strings.ReplaceAll(lower, "&", " and ")
	lower = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(lower, "-")
	lower = strings.Trim(lower, "-")
	if lower == "" {
		return "organizer"
	}
	return lower
}

func (o *Organizer) UpdateProfile(name, phone, slug string) error {
	name = strings.TrimSpace(name)
	phone = strings.TrimSpace(phone)
	if slug != "" {
		slug = strings.TrimSpace(slug)
	}
	if name == "" {
		return ErrNameRequired
	}
	if slug == "" {
		slug = toSlug(name)
	}
	if err := ValidateSlug(slug); err != nil {
		return err
	}

	o.Name = name
	o.Phone = phone
	o.Slug = slug
	o.UpdatedAt = time.Now().UTC()
	return nil
}

func (o *Organizer) Activate() error {
	if o.IsActive {
		return ErrOrganizerAlreadyActive
	}
	o.IsActive = true
	o.UpdatedAt = time.Now().UTC()
	return nil
}

func (o *Organizer) Deactivate() error {
	if !o.IsActive {
		return ErrOrganizerAlreadyInactive
	}
	o.IsActive = false
	o.UpdatedAt = time.Now().UTC()
	return nil
}
