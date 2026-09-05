package user

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")

	ErrEmailRequired = errors.New(
		"email is required",
	)

	ErrInvalidEmail = errors.New(
		"invalid email",
	)

	ErrPasswordHashRequired = errors.New(
		"password hash is required",
	)

	ErrFullNameRequired = errors.New(
		"full name is required",
	)

	ErrUserInactive = errors.New(
		"user is inactive",
	)

	ErrUserAlreadyActive = errors.New(
		"user is already active",
	)

	ErrUserAlreadyInactive = errors.New(
		"user is already inactive",
	)
)
