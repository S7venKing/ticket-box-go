package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
)

type EventRepository struct {
	db *sql.DB
}

var _ event.Repository = (*EventRepository)(nil)

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

const selectEventColumns = `
	SELECT
		id,
		organizer_id,
		title,
		description,
		venue,
		start_at,
		end_at,
		capacity,
		status,
		created_at,
		updated_at
	FROM events
`

func (r *EventRepository) Create(ctx context.Context, e *event.Event) error {
	const query = `
		INSERT INTO events (
			id,
			organizer_id,
			title,
			description,
			venue,
			start_at,
			end_at,
			capacity,
			status,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		e.ID.String(),
		e.OrganizerID.String(),
		e.Title,
		e.Description,
		e.Venue,
		e.StartAt,
		e.EndAt,
		e.Capacity,
		string(e.Status),
		e.CreatedAt,
		e.UpdatedAt,
	)
	return err
}

func (r *EventRepository) GetByID(ctx context.Context, id uuid.UUID) (*event.Event, error) {
	row := r.db.QueryRowContext(ctx, selectEventColumns+` WHERE id = ? LIMIT 1`, id.String())
	return scanEvent(row)
}

func (r *EventRepository) GetByOrganizer(ctx context.Context, organizerID uuid.UUID, offset, limit int) ([]*event.Event, error) {
	if limit <= 0 {
		limit = 20
	}
	query := selectEventColumns + ` WHERE organizer_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, organizerID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*event.Event
	for rows.Next() {
		e, err := scanEventFromRows(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *EventRepository) Update(ctx context.Context, e *event.Event) error {
	const query = `
		UPDATE events
		SET
			title = ?,
			description = ?,
			venue = ?,
			start_at = ?,
			end_at = ?,
			capacity = ?,
			status = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		e.Title,
		e.Description,
		e.Venue,
		e.StartAt,
		e.EndAt,
		e.Capacity,
		string(e.Status),
		e.UpdatedAt,
		e.ID.String(),
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return event.ErrEventNotFound
	}
	return nil
}

func (r *EventRepository) List(ctx context.Context, offset, limit int) ([]*event.Event, error) {
	if limit <= 0 {
		limit = 20
	}
	query := selectEventColumns + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*event.Event
	for rows.Next() {
		e, err := scanEventFromRows(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func scanEvent(row *sql.Row) (*event.Event, error) {
	var (
		e         event.Event
		idString  string
		orgID     string
		title     string
		desc      sql.NullString
		venue     string
		startAt   time.Time
		endAt     time.Time
		capacity  int
		status    string
		createdAt time.Time
		updatedAt time.Time
	)

	err := row.Scan(
		&idString,
		&orgID,
		&title,
		&desc,
		&venue,
		&startAt,
		&endAt,
		&capacity,
		&status,
		&createdAt,
		&updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, event.ErrEventNotFound
	}
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return nil, err
	}
	organizerID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, err
	}

	e = event.Event{
		ID:          id,
		OrganizerID: organizerID,
		Title:       title,
		Description: desc.String,
		Venue:       venue,
		StartAt:     startAt,
		EndAt:       endAt,
		Capacity:    capacity,
		Status:      event.Status(status),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	return &e, nil
}

func scanEventFromRows(rows *sql.Rows) (*event.Event, error) {
	var (
		e         event.Event
		idString  string
		orgID     string
		title     string
		desc      sql.NullString
		venue     string
		startAt   time.Time
		endAt     time.Time
		capacity  int
		status    string
		createdAt time.Time
		updatedAt time.Time
	)

	err := rows.Scan(
		&idString,
		&orgID,
		&title,
		&desc,
		&venue,
		&startAt,
		&endAt,
		&capacity,
		&status,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return nil, err
	}
	organizerID, err := uuid.Parse(orgID)
	if err != nil {
		return nil, err
	}

	e = event.Event{
		ID:          id,
		OrganizerID: organizerID,
		Title:       title,
		Description: desc.String,
		Venue:       venue,
		StartAt:     startAt,
		EndAt:       endAt,
		Capacity:    capacity,
		Status:      event.Status(status),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	return &e, nil
}
