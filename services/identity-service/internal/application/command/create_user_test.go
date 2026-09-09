package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/command"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
)

// fakeRepo embeds the interface so only the methods CreateUser touches are implemented.
type fakeRepo struct {
	user.Repository

	exists  bool
	created *user.User
}

func (f *fakeRepo) ExistsByEmail(context.Context, string) (bool, error) { return f.exists, nil }

func (f *fakeRepo) Create(_ context.Context, u *user.User) error {
	f.created = u
	return nil
}

type fakeHasher struct{}

func (fakeHasher) Hash(password string) (string, error) { return "hashed:" + password, nil }
func (fakeHasher) Compare(string, string) error         { return nil }

func TestCreateUser_PasswordValidation(t *testing.T) {
	h := command.NewCreateUserHandler(&fakeRepo{}, fakeHasher{})

	_, err := h.Handle(context.Background(), command.CreateUserCommand{Email: "a@b.com", FullName: "A"})
	require.ErrorIs(t, err, command.ErrPasswordRequired)

	_, err = h.Handle(context.Background(), command.CreateUserCommand{Email: "a@b.com", Password: "short", FullName: "A"})
	require.ErrorIs(t, err, command.ErrPasswordTooShort)
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	h := command.NewCreateUserHandler(&fakeRepo{exists: true}, fakeHasher{})

	_, err := h.Handle(context.Background(), command.CreateUserCommand{
		Email: "a@b.com", Password: "correct-horse", FullName: "A",
	})
	require.ErrorIs(t, err, user.ErrEmailAlreadyExists)
}

func TestCreateUser_Success(t *testing.T) {
	repo := &fakeRepo{}
	h := command.NewCreateUserHandler(repo, fakeHasher{})

	result, err := h.Handle(context.Background(), command.CreateUserCommand{
		Email: " User@Example.com ", Password: "correct-horse", FullName: "Nguyen Van A", Phone: "0900000000",
	})
	require.NoError(t, err)

	require.Equal(t, "user@example.com", result.Email)
	require.True(t, result.IsActive)
	require.NotNil(t, repo.created)
	require.Equal(t, "hashed:correct-horse", repo.created.PasswordHash)
}
