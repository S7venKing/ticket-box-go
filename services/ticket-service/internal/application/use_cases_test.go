package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/command"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/query"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/event"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/domain/organizer"
)

type fakeOrganizerRepository struct {
	items map[string]*organizer.Organizer
}

func newFakeOrganizerRepository() *fakeOrganizerRepository {
	return &fakeOrganizerRepository{items: map[string]*organizer.Organizer{}}
}

func (r *fakeOrganizerRepository) Create(ctx context.Context, o *organizer.Organizer) error {
	r.items[o.ID.String()] = o
	return nil
}

func (r *fakeOrganizerRepository) GetByID(ctx context.Context, id uuid.UUID) (*organizer.Organizer, error) {
	o, ok := r.items[id.String()]
	if !ok {
		return nil, organizer.ErrOrganizerNotFound
	}
	return o, nil
}

func (r *fakeOrganizerRepository) GetByEmail(ctx context.Context, email string) (*organizer.Organizer, error) {
	for _, o := range r.items {
		if o.Email == organizer.NormalizeEmail(email) {
			return o, nil
		}
	}
	return nil, organizer.ErrOrganizerNotFound
}

func (r *fakeOrganizerRepository) Update(ctx context.Context, o *organizer.Organizer) error {
	r.items[o.ID.String()] = o
	return nil
}

func (r *fakeOrganizerRepository) List(ctx context.Context, offset, limit int) ([]*organizer.Organizer, error) {
	result := make([]*organizer.Organizer, 0, len(r.items))
	for _, o := range r.items {
		result = append(result, o)
	}
	return result, nil
}

type fakeEventRepository struct {
	items map[string]*event.Event
}

func newFakeEventRepository() *fakeEventRepository {
	return &fakeEventRepository{items: map[string]*event.Event{}}
}

func (r *fakeEventRepository) Create(ctx context.Context, e *event.Event) error {
	r.items[e.ID.String()] = e
	return nil
}

func (r *fakeEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*event.Event, error) {
	e, ok := r.items[id.String()]
	if !ok {
		return nil, event.ErrEventNotFound
	}
	return e, nil
}

func (r *fakeEventRepository) GetByOrganizer(ctx context.Context, organizerID uuid.UUID, offset, limit int) ([]*event.Event, error) {
	result := make([]*event.Event, 0)
	for _, e := range r.items {
		if e.OrganizerID == organizerID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (r *fakeEventRepository) Update(ctx context.Context, e *event.Event) error {
	r.items[e.ID.String()] = e
	return nil
}

func (r *fakeEventRepository) List(ctx context.Context, offset, limit int) ([]*event.Event, error) {
	result := make([]*event.Event, 0, len(r.items))
	for _, e := range r.items {
		result = append(result, e)
	}
	return result, nil
}

func TestOrganizerLifecycleHandlers(t *testing.T) {
	repo := newFakeOrganizerRepository()
	o, err := organizer.NewOrganizer("Tech Meetup", "team@example.com", "0900000001", "tech-meetup")
	require.NoError(t, err)
	require.NoError(t, repo.Create(context.Background(), o))

	deactivate := command.NewDeactivateOrganizerHandler(repo)
	updated, err := deactivate.Handle(context.Background(), command.DeactivateOrganizerCommand{ID: o.ID})
	require.NoError(t, err)
	require.False(t, updated.IsActive)

	activate := command.NewActivateOrganizerHandler(repo)
	updated, err = activate.Handle(context.Background(), command.ActivateOrganizerCommand{ID: o.ID})
	require.NoError(t, err)
	require.True(t, updated.IsActive)

	byEmail := query.NewGetOrganizerByEmailHandler(repo)
	found, err := byEmail.Handle(context.Background(), query.GetOrganizerByEmailQuery{Email: "team@example.com"})
	require.NoError(t, err)
	require.Equal(t, o.ID, found.ID)
}

func TestEventLifecycleAndQueryHandlers(t *testing.T) {
	repo := newFakeEventRepository()
	e, err := event.NewEvent(uuid.New(), "Launch Night", "Welcome", "HCM", time.Now().Add(2*time.Hour), time.Now().Add(3*time.Hour), 120)
	require.NoError(t, err)
	require.NoError(t, repo.Create(context.Background(), e))

	update := command.NewUpdateEventDetailsHandler(repo)
	updated, err := update.Handle(context.Background(), command.UpdateEventDetailsCommand{
		ID:          e.ID,
		Title:       "Launch Night 2",
		Description: "Updated",
		Venue:       "Da Nang",
		StartAt:     time.Now().Add(10 * time.Hour),
		EndAt:       time.Now().Add(12 * time.Hour),
		Capacity:    200,
	})
	require.NoError(t, err)
	require.Equal(t, "Launch Night 2", updated.Title)
	require.Equal(t, 200, updated.Capacity)

	cancel := command.NewCancelEventHandler(repo)
	updated, err = cancel.Handle(context.Background(), command.CancelEventCommand{ID: e.ID})
	require.NoError(t, err)
	require.Equal(t, event.StatusCancelled, updated.Status)

	close := command.NewCloseEventHandler(repo)
	updated, err = close.Handle(context.Background(), command.CloseEventCommand{ID: e.ID})
	require.NoError(t, err)
	require.Equal(t, event.StatusClosed, updated.Status)

	list := query.NewListEventsHandler(repo)
	items, err := list.Handle(context.Background(), query.ListEventsQuery{Offset: 0, Limit: 10})
	require.NoError(t, err)
	require.Len(t, items, 1)

	byOrganizer := query.NewGetEventsByOrganizerHandler(repo)
	organizerEvents, err := byOrganizer.Handle(context.Background(), query.GetEventsByOrganizerQuery{OrganizerID: e.OrganizerID, Offset: 0, Limit: 10})
	require.NoError(t, err)
	require.Len(t, organizerEvents, 1)
}
