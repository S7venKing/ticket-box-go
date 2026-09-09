package grpc

import (
	"errors"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, user.ErrUserNotFound):
		return status.Error(
			codes.NotFound,
			"user not found",
		)

	case errors.Is(err, user.ErrEmailRequired):
		return status.Error(
			codes.InvalidArgument,
			"email is required",
		)

	case errors.Is(err, user.ErrInvalidEmail):
		return status.Error(
			codes.InvalidArgument,
			"invalid email",
		)

	case errors.Is(err, user.ErrFullNameRequired):
		return status.Error(
			codes.InvalidArgument,
			"full name is required",
		)

	case errors.Is(err, user.ErrPasswordHashRequired):
		return status.Error(
			codes.Internal,
			"password configuration error",
		)

	case errors.Is(err, user.ErrUserAlreadyActive):
		return status.Error(
			codes.FailedPrecondition,
			"user is already active",
		)

	case errors.Is(err, user.ErrUserAlreadyInactive):
		return status.Error(
			codes.FailedPrecondition,
			"user is already inactive",
		)

	default:
		return status.Error(
			codes.Internal,
			"internal server error",
		)
	}
}
