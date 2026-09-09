package user_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
)

func TestNewUser_NormalizesEmail(t *testing.T) {
	u, err := user.NewUser("  User@Example.COM ", "hashed-password", "Nguyen Van A", "")

	require.NoError(t, err)
	require.Equal(t, "user@example.com", u.Email)
	require.True(t, u.IsActive)
}

func TestNewUser_InvalidEmail(t *testing.T) {
	_, err := user.NewUser("invalid-email", "hashed-password", "Nguyen Van A", "")
	require.ErrorIs(t, err, user.ErrInvalidEmail)

	_, err = user.NewUser("Name <user@example.com>", "hashed-password", "Nguyen Van A", "")
	require.ErrorIs(t, err, user.ErrInvalidEmail)

	_, err = user.NewUser("", "hashed-password", "Nguyen Van A", "")
	require.ErrorIs(t, err, user.ErrEmailRequired)
}

func TestNewUser_EmptyFullName(t *testing.T) {
	_, err := user.NewUser("user@example.com", "hashed-password", "", "")

	require.ErrorIs(t, err, user.ErrFullNameRequired)
}

func TestNewUser_EmptyPasswordHash(t *testing.T) {
	_, err := user.NewUser("user@example.com", "", "Nguyen Van A", "")

	require.ErrorIs(t, err, user.ErrPasswordHashRequired)
}

func TestUser_Deactivate(t *testing.T) {
	u, err := user.NewUser("user@example.com", "hashed-password", "Nguyen Van A", "")
	require.NoError(t, err)
	require.True(t, u.IsActive)

	require.NoError(t, u.Deactivate())
	require.False(t, u.IsActive)

	require.ErrorIs(t, u.Deactivate(), user.ErrUserAlreadyInactive)
}

func TestUser_Activate_AlreadyActive(t *testing.T) {
	u, err := user.NewUser("user@example.com", "hashed-password", "Nguyen Van A", "")
	require.NoError(t, err)

	require.ErrorIs(t, u.Activate(), user.ErrUserAlreadyActive)
}
