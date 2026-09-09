package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

var _ user.Repository = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	u *user.User,
) error {
	const query = `
		INSERT INTO users (
			id,
			email,
			password_hash,
			full_name,
			phone,
			is_active,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		u.ID.String(),
		u.Email,
		u.PasswordHash,
		u.FullName,
		u.Phone,
		u.IsActive,
		u.CreatedAt,
		u.UpdatedAt,
	)

	return err
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*user.User, error) {
	const query = `
		SELECT
			id,
			email,
			password_hash,
			full_name,
			phone,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE id = ?
		LIMIT 1
	`

	var (
		u         user.User
		idString  string
		phone     sql.NullString
		isActive  bool
		createdAt = u.CreatedAt
		updatedAt = u.UpdatedAt
	)

	err := r.db.QueryRowContext(
		ctx,
		query,
		id.String(),
	).Scan(
		&idString,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&phone,
		&isActive,
		&createdAt,
		&updatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	parsedID, err := uuid.Parse(idString)
	if err != nil {
		return nil, err
	}

	u.ID = parsedID
	u.Phone = phone.String
	u.IsActive = isActive
	u.CreatedAt = createdAt
	u.UpdatedAt = updatedAt

	return &u, nil
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*user.User, error) {
	const query = `
		SELECT
			id,
			email,
			password_hash,
			full_name,
			phone,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	var (
		u        user.User
		idString string
		phone    sql.NullString
	)

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&idString,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&phone,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return nil, err
	}

	u.ID = id
	u.Phone = phone.String

	return &u, nil
}

func (r *UserRepository) Update(
	ctx context.Context,
	u *user.User,
) error {
	const query = `
		UPDATE users
		SET
			email = ?,
			password_hash = ?,
			full_name = ?,
			phone = ?,
			is_active = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		u.Email,
		u.PasswordHash,
		u.FullName,
		u.Phone,
		u.IsActive,
		u.UpdatedAt,
		u.ID.String(),
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) ExistsByEmail(
	ctx context.Context,
	email string,
) (bool, error) {
	const query = `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email = ?
		)
	`

	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(&exists)

	return exists, err
}
