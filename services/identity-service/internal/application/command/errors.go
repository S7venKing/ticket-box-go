package command

import "errors"

// MinPasswordLength is the minimum accepted plaintext password length
// (NIST SP 800-63B recommends at least 8 characters).
const MinPasswordLength = 8

var (
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password is too short")
)
