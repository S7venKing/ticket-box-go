package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/command"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
)

func TestMapError(t *testing.T) {
	s := &IdentityServer{logger: slog.New(slog.NewTextHandler(io.Discard, nil))}

	sqlErr := fmt.Errorf("query users: %w", errors.New("Error 1045: Access denied for user 'root'@'10.0.0.1'"))

	cases := []struct {
		name string
		err  error
		code codes.Code
	}{
		{"not found", user.ErrUserNotFound, codes.NotFound},
		{"email exists", user.ErrEmailAlreadyExists, codes.AlreadyExists},
		{"invalid email", user.ErrInvalidEmail, codes.InvalidArgument},
		{"password too short", command.ErrPasswordTooShort, codes.InvalidArgument},
		{"already active", user.ErrUserAlreadyActive, codes.FailedPrecondition},
		{"deadline", context.DeadlineExceeded, codes.DeadlineExceeded},
		{"wrapped sql error", sqlErr, codes.Internal},
		{"password hash invariant", user.ErrPasswordHashRequired, codes.Internal},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			st, ok := status.FromError(s.mapError(context.Background(), tc.err))
			require.True(t, ok)
			require.Equal(t, tc.code, st.Code())

			if tc.code == codes.Internal {
				require.Equal(t, "internal server error", st.Message())
			}
		})
	}
}
