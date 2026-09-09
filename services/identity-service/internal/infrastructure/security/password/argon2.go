package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/command"
	"golang.org/x/crypto/argon2"
)

const (
	memory      = 64 * 1024
	iterations  = 3
	parallelism = 2
	saltLength  = 16
	keyLength   = 32
)

type Argon2Hasher struct{}

var _ command.PasswordHasher = (*Argon2Hasher)(nil)

func NewArgon2Hasher() *Argon2Hasher {
	return &Argon2Hasher{}
}

func (h *Argon2Hasher) Hash(
	password string,
) (string, error) {
	if password == "" {
		return "", errors.New("password is required")
	}

	salt := make([]byte, saltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLength,
	)

	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)
	hashEncoded := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		memory,
		iterations,
		parallelism,
		saltEncoded,
		hashEncoded,
	), nil
}

func (h *Argon2Hasher) Compare(
	password string,
	encodedHash string,
) error {
	if password == "" {
		return errors.New("password is required")
	}

	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return errors.New("invalid password hash")
	}

	if parts[1] != "argon2id" {
		return errors.New("invalid password hash algorithm")
	}

	if parts[2] != "v=19" {
		return errors.New("unsupported argon2 version")
	}

	params := strings.Split(parts[3], ",")

	if len(params) != 3 {
		return errors.New("invalid argon2 parameters")
	}

	memoryValue, err := parseArgon2Parameter(params[0], "m")
	if err != nil {
		return err
	}

	iterationsValue, err := parseArgon2Parameter(params[1], "t")
	if err != nil {
		return err
	}

	parallelismValue, err := parseArgon2Parameter(params[2], "p")
	if err != nil {
		return err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return errors.New("invalid password salt")
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return errors.New("invalid password hash")
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		uint32(iterationsValue),
		uint32(memoryValue),
		uint8(parallelismValue),
		uint32(len(expectedHash)),
	)

	if subtle.ConstantTimeCompare(
		actualHash,
		expectedHash,
	) != 1 {
		return errors.New("invalid password")
	}

	return nil
}

func parseArgon2Parameter(
	value string,
	prefix string,
) (int, error) {
	if !strings.HasPrefix(value, prefix+"=") {
		return 0, errors.New("invalid argon2 parameter")
	}

	raw := strings.TrimPrefix(
		value,
		prefix+"=",
	)

	result, err := strconv.Atoi(raw)
	if err != nil {
		return 0, errors.New("invalid argon2 parameter")
	}

	if result <= 0 {
		return 0, errors.New("invalid argon2 parameter")
	}

	return result, nil
}
