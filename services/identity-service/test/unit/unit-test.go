package unit

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/s7venking/ticket-box/identity/internal/domain/user"
)

func TestNewUser_InvalidEmail(t *testing.T) {
	_, err := user.NewUser(
		"invalid-email",
		"hashed-password",
		"Nguyen Van A",
		"",
	)

	require.ErrorIs(
		t,
		err,
		user.ErrInvalidEmail,
	)
}

func TestNewUser_EmptyFullName(t *testing.T) {
	_, err := user.NewUser(
		"user@example.com",
		"hashed-password",
		"",
		"",
	)

	require.ErrorIs(
		t,
		err,
		user.ErrFullNameRequired,
	)
}

func TestUser_Deactivate(t *testing.T) {
	u, err := user.NewUser(
		"user@example.com",
		"hashed-password",
		"Nguyen Van A",
		"",
	)

	require.NoError(t, err)
	require.True(t, u.IsActive)

	err = u.Deactivate()

	require.NoError(t, err)
	require.False(t, u.IsActive)
}

func TestUser_Deactivate_AlreadyInactive(t *testing.T) {
	u, err := user.NewUser(
		"user@example.com",
		"hashed-password",
		"Nguyen Van A",
		"",
	)

	require.NoError(t, err)

	require.NoError(t, u.Deactivate())

	err = u.Deactivate()

	require.ErrorIs(
		t,
		err,
		user.ErrUserAlreadyInactive,
	)
}
