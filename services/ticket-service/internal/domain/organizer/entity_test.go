package organizer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/organizer"
)

func TestNewOrganizer_Valid(t *testing.T) {
	o, err := organizer.NewOrganizer("Tech Meetup VN", "  TEAM@EXAMPLE.COM  ", "0900000001", "")
	require.NoError(t, err)
	require.Equal(t, "team@example.com", o.Email)
	require.Equal(t, "tech-meetup-vn", o.Slug)
	require.True(t, o.IsActive)
	require.False(t, o.CreatedAt.IsZero())
	require.False(t, o.UpdatedAt.IsZero())
}

func TestNewOrganizer_InvalidEmail(t *testing.T) {
	_, err := organizer.NewOrganizer("Name", "bad-email", "0900000001", "")
	require.ErrorIs(t, err, organizer.ErrInvalidEmail)
}

func TestOrganizer_UpdateProfile(t *testing.T) {
	o, err := organizer.NewOrganizer("Old Name", "team@example.com", "0900000001", "old-name")
	require.NoError(t, err)

	before := o.UpdatedAt
	time.Sleep(10 * time.Millisecond)
	require.NoError(t, o.UpdateProfile("New Name", "0909999999", "new-name"))
	require.Equal(t, "New Name", o.Name)
	require.Equal(t, "new-name", o.Slug)
	require.True(t, o.UpdatedAt.After(before))
}

func TestOrganizer_Deactivate(t *testing.T) {
	o, err := organizer.NewOrganizer("Name", "team@example.com", "0900000001", "name")
	require.NoError(t, err)
	require.NoError(t, o.Deactivate())
	require.False(t, o.IsActive)
	require.ErrorIs(t, o.Deactivate(), organizer.ErrOrganizerAlreadyInactive)
}
