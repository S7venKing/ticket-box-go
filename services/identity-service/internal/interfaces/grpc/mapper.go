package grpc

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/command"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/organizer"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
)

// mapError translates domain/application errors into gRPC statuses.
//
// Anything not explicitly known (SQL/driver errors, hasher failures, bugs) is
// logged server-side WITH its cause and returned to the client as
// codes.Internal with a fixed message, so implementation details never leak.
func (s *IdentityServer) mapError(ctx context.Context, err error) error {
	if code, msg, ok := knownError(err); ok {
		return status.Error(code, msg)
	}

	method, _ := grpc.Method(ctx)

	s.logger.ErrorContext(ctx, "unhandled error",
		slog.String("method", method),
		slog.Any("error", err),
	)

	return status.Error(codes.Internal, "internal server error")
}

func knownError(err error) (codes.Code, string, bool) {
	switch {
	case errors.Is(err, user.ErrUserNotFound):
		return codes.NotFound, "user not found", true

	case errors.Is(err, user.ErrEmailAlreadyExists):
		return codes.AlreadyExists, "email already exists", true
	case errors.Is(err, organizer.ErrOrganizerNotFound):
		return codes.NotFound, "organizer not found", true
	case errors.Is(err, organizer.ErrOrganizerDuplicate):
		return codes.AlreadyExists, "organizer email or slug already exists", true
	case errors.Is(err, organizer.ErrNameRequired):
		return codes.InvalidArgument, "organizer name is required", true
	case errors.Is(err, organizer.ErrEmailRequired):
		return codes.InvalidArgument, "organizer email is required", true
	case errors.Is(err, organizer.ErrInvalidEmail):
		return codes.InvalidArgument, "organizer email is invalid", true
	case errors.Is(err, organizer.ErrSlugRequired):
		return codes.InvalidArgument, "organizer slug is required", true
	case errors.Is(err, organizer.ErrInvalidSlug):
		return codes.InvalidArgument, "organizer slug is invalid", true

	case errors.Is(err, user.ErrEmailRequired):
		return codes.InvalidArgument, "email is required", true
	case errors.Is(err, user.ErrInvalidEmail):
		return codes.InvalidArgument, "invalid email", true
	case errors.Is(err, user.ErrFullNameRequired):
		return codes.InvalidArgument, "full name is required", true
	case errors.Is(err, command.ErrPasswordRequired):
		return codes.InvalidArgument, "password is required", true
	case errors.Is(err, command.ErrPasswordTooShort):
		return codes.InvalidArgument, "password is too short", true
	case errors.Is(err, command.ErrInvalidCredentials):
		return codes.Unauthenticated, "invalid credentials", true

	case errors.Is(err, user.ErrUserAlreadyActive):
		return codes.FailedPrecondition, "user is already active", true
	case errors.Is(err, user.ErrUserAlreadyInactive):
		return codes.FailedPrecondition, "user is already inactive", true
	case errors.Is(err, user.ErrUserInactive):
		return codes.FailedPrecondition, "user is inactive", true

	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded, "deadline exceeded", true
	case errors.Is(err, context.Canceled):
		return codes.Canceled, "request canceled", true
	}

	return codes.Unknown, "", false
}
