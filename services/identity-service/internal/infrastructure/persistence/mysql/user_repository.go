package mysql

import (
	"context"
	"database/sql"
	"errors"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
)

// mysqlErrDuplicateEntry is MySQL error 1062 (ER_DUP_ENTRY).
const mysqlErrDuplicateEntry = 1062

type UserRepository struct {
	db *sql.DB
}

var _ user.Repository = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

const selectUserColumns = `
	SELECT
		id,
		email,
		password_hash,
		full_name,
		phone,
		is_active,
		role,
		organizer_id,
		created_at,
		updated_at
	FROM users
`

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
			role,
			organizer_id,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		u.Role,
		nullableUUID(u.OrganizerID),
		u.CreatedAt,
		u.UpdatedAt,
	)

	if isDuplicateEntry(err) {
		return user.ErrEmailAlreadyExists
	}

	return err
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*user.User, error) {
	row := r.db.QueryRowContext(ctx, selectUserColumns+` WHERE id = ? LIMIT 1`, id.String())

	return scanUser(row)
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*user.User, error) {
	row := r.db.QueryRowContext(ctx, selectUserColumns+` WHERE email = ? LIMIT 1`, email)

	return scanUser(row)
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	rows, err := r.db.QueryContext(ctx, selectUserColumns+` ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*user.User, 0)
	for rows.Next() {
		item, err := scanUserRows(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
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
			,organizer_id = ?
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
		nullableUUID(u.OrganizerID),
		u.ID.String(),
	)

	if isDuplicateEntry(err) {
		return user.ErrEmailAlreadyExists
	}

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
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`

	var exists bool

	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)

	return exists, err
}

func scanUser(row *sql.Row) (*user.User, error) {
	var (
		u           user.User
		idString    string
		phone       sql.NullString
		organizerID sql.NullString
	)

	err := row.Scan(
		&idString,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&phone,
		&u.IsActive,
		&u.Role,
		&organizerID,
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
	u.OrganizerID = parseNullableUUID(organizerID)

	return &u, nil
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUserRows(row userScanner) (*user.User, error) {
	var (
		u           user.User
		idString    string
		phone       sql.NullString
		organizerID sql.NullString
	)
	if err := row.Scan(
		&idString,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&phone,
		&u.IsActive,
		&u.Role,
		&organizerID,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		return nil, err
	}
	u.ID = id
	u.OrganizerID = parseNullableUUID(organizerID)
	if phone.Valid {
		u.Phone = phone.String
	}

	return &u, nil
}

func nullableUUID(id *uuid.UUID) any {
	if id == nil {
		return nil
	}
	return id.String()
}

func parseNullableUUID(value sql.NullString) *uuid.UUID {
	if !value.Valid || value.String == "" {
		return nil
	}
	id, err := uuid.Parse(value.String)
	if err != nil {
		return nil
	}
	return &id
}

// isDuplicateEntry translates the driver-specific unique-violation into a
// signal the domain understands. This is the only place that inspects
// MySQL error codes.
func isDuplicateEntry(err error) bool {
	var mysqlErr *mysqldriver.MySQLError

	return errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlErrDuplicateEntry
}
