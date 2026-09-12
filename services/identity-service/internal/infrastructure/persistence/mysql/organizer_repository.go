package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/organizer"
	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type OrganizerRepository struct{ db *sql.DB }

var _ organizer.Repository = (*OrganizerRepository)(nil)

func NewOrganizerRepository(db *sql.DB) *OrganizerRepository { return &OrganizerRepository{db: db} }

const selectOrganizerColumns = `SELECT id, email, name, phone, slug, is_active, created_at, updated_at FROM organizers`

func (r *OrganizerRepository) Create(ctx context.Context, o *organizer.Organizer) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO organizers (id,email,name,phone,slug,is_active,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?)`, o.ID.String(), o.Email, o.Name, o.Phone, o.Slug, o.IsActive, o.CreatedAt, o.UpdatedAt)
	if isOrganizerDuplicate(err) {
		return organizer.ErrOrganizerDuplicate
	}
	return err
}
func (r *OrganizerRepository) GetByID(ctx context.Context, id uuid.UUID) (*organizer.Organizer, error) {
	return scanIdentityOrganizer(r.db.QueryRowContext(ctx, selectOrganizerColumns+` WHERE id=? LIMIT 1`, id.String()))
}
func (r *OrganizerRepository) Update(ctx context.Context, o *organizer.Organizer) error {
	result, err := r.db.ExecContext(ctx, `UPDATE organizers SET email=?,name=?,phone=?,slug=?,is_active=?,updated_at=? WHERE id=?`, o.Email, o.Name, o.Phone, o.Slug, o.IsActive, o.UpdatedAt, o.ID.String())
	if isOrganizerDuplicate(err) {
		return organizer.ErrOrganizerDuplicate
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
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 100
	}
	rows, err := r.db.QueryContext(ctx, selectOrganizerColumns+` ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]*organizer.Organizer, 0)
	for rows.Next() {
		o, err := scanIdentityOrganizerValues(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, o)
	}
	return items, rows.Err()
}

type organizerScanner interface{ Scan(...any) error }

func scanIdentityOrganizer(row *sql.Row) (*organizer.Organizer, error) {
	return scanIdentityOrganizerValues(row)
}
func scanIdentityOrganizerValues(scanner organizerScanner) (*organizer.Organizer, error) {
	var o organizer.Organizer
	var id string
	var phone sql.NullString
	if err := scanner.Scan(&id, &o.Email, &o.Name, &phone, &o.Slug, &o.IsActive, &o.CreatedAt, &o.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, organizer.ErrOrganizerNotFound
		}
		return nil, err
	}
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	o.ID, o.Phone = parsed, phone.String
	return &o, nil
}
func isOrganizerDuplicate(err error) bool {
	var mysqlErr *mysqldriver.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
