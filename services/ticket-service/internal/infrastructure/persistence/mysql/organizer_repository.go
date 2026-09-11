package mysql

import (
	"context"
	"database/sql"
	"errors"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/organizer"
)

const mysqlErrDuplicateEntry = 1062

type OrganizerRepository struct {
	db *sql.DB
}

var _ organizer.Repository = (*OrganizerRepository)(nil)

func NewOrganizerRepository(db *sql.DB) *OrganizerRepository {
	return &OrganizerRepository{db: db}
}

const selectOrganizerColumns = `
	SELECT
		id,
		email,
		name,
		phone,
		slug,
		is_active,
		created_at,
		updated_at
	FROM organizers
`

func (r *OrganizerRepository) Create(ctx context.Context, o *organizer.Organizer) error {
	const query = `
		INSERT INTO organizers (
			id,
			email,
			name,
			phone,
			slug,
			is_active,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		o.ID.String(),
		o.Email,
		o.Name,
		o.Phone,
		o.Slug,
		o.IsActive,
		o.CreatedAt,
		o.UpdatedAt,
	)
	if isDuplicateEntry(err) {
		return organizer.ErrEmailAlreadyExists
	}
	return err
}

func (r *OrganizerRepository) GetByID(ctx context.Context, id uuid.UUID) (*organizer.Organizer, error) {
	row := r.db.QueryRowContext(ctx, selectOrganizerColumns+` WHERE id = ? LIMIT 1`, id.String())
	return scanOrganizer(row)
}

func (r *OrganizerRepository) GetByEmail(ctx context.Context, email string) (*organizer.Organizer, error) {
	row := r.db.QueryRowContext(ctx, selectOrganizerColumns+` WHERE email = ? LIMIT 1`, email)
	return scanOrganizer(row)
}

func (r *OrganizerRepository) Update(ctx context.Context, o *organizer.Organizer) error {
	const query = `
		UPDATE organizers
		SET
			email = ?,
			name = ?,
			phone = ?,
			slug = ?,
			is_active = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		o.Email,
		o.Name,
		o.Phone,
		o.Slug,
		o.IsActive,
		o.UpdatedAt,
		o.ID.String(),
	)
	if isDuplicateEntry(err) {
		return organizer.ErrEmailAlreadyExists
	}
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return organizer.ErrOrganizerNotFound
	}
	return nil
}

func (r *OrganizerRepository) List(ctx context.Context, offset, limit int) ([]*organizer.Organizer, error) {
	if limit <= 0 {
		limit = 20
	}
	query := selectOrganizerColumns + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*organizer.Organizer
	for rows.Next() {
		o, err := scanOrganizerFromRows(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, o)
	}
	return result, rows.Err()
}

func scanOrganizer(row *sql.Row) (*organizer.Organizer, error) {
	var (
		o        organizer.Organizer
		idString string
		phone    sql.NullString
	)

	err := row.Scan(
		&idString,
		&o.Email,
		&o.Name,
		&phone,
		&o.Slug,
		&o.IsActive,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, organizer.ErrOrganizerNotFound
	}
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		return nil, err
	}
	o.ID = id
	o.Phone = phone.String
	return &o, nil
}

func scanOrganizerFromRows(rows *sql.Rows) (*organizer.Organizer, error) {
	var (
		o        organizer.Organizer
		idString string
		phone    sql.NullString
	)

	err := rows.Scan(
		&idString,
		&o.Email,
		&o.Name,
		&phone,
		&o.Slug,
		&o.IsActive,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	id, err := uuid.Parse(idString)
	if err != nil {
		return nil, err
	}
	o.ID = id
	o.Phone = phone.String
	return &o, nil
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysqldriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlErrDuplicateEntry
}
