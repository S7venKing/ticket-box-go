package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	identityv1 "github.com/S7venKing/ticket-box-go/gen/identity/v1"
	ticketv1 "github.com/S7venKing/ticket-box-go/gen/ticket/v1"
	"github.com/S7venKing/ticket-box-go/services/api-gateway/internal/auth"
)

type fakeIdentityClient struct {
	createUserFn func(context.Context, string, string, string, string, string) (*identityv1.User, error)
	authFn       func(context.Context, string, string) (*identityv1.User, string, error)
	getUserFn    func(context.Context, string) (*identityv1.User, error)
}

func (f *fakeIdentityClient) CreateUser(ctx context.Context, email, password, fullName, phone, role string) (*identityv1.User, error) {
	if f.createUserFn != nil {
		return f.createUserFn(ctx, email, password, fullName, phone, role)
	}
	return &identityv1.User{Id: "user-1", Email: email, FullName: fullName, Phone: phone, IsActive: true}, nil
}

func (f *fakeIdentityClient) AuthenticateUser(ctx context.Context, email, password string) (*identityv1.User, string, error) {
	if f.authFn != nil {
		return f.authFn(ctx, email, password)
	}
	return &identityv1.User{Id: "user-1", Email: email}, "identity-token", nil
}

func (f *fakeIdentityClient) GetUserByID(ctx context.Context, id string) (*identityv1.User, error) {
	if f.getUserFn != nil {
		return f.getUserFn(ctx, id)
	}
	return &identityv1.User{Id: id, Email: "user@example.com"}, nil
}

type fakeTicketClient struct {
	createOrganizerFn  func(context.Context, string, string, string, string) (*ticketv1.Organizer, error)
	getOrganizerByIDFn func(context.Context, string) (*ticketv1.Organizer, error)
	listOrganizersFn   func(context.Context, int32, int32) ([]*ticketv1.Organizer, error)
	updateOrganizerFn  func(context.Context, string, string, string, string) (*ticketv1.Organizer, error)
	createEventFn      func(context.Context, string, string, string, string, time.Time, time.Time, int32) (*ticketv1.Event, error)
	getEventByIDFn     func(context.Context, string) (*ticketv1.Event, error)
	listEventsFn       func(context.Context, int32, int32) ([]*ticketv1.Event, error)
	publishEventFn     func(context.Context, string) (*ticketv1.Event, error)
	submitEventFn      func(context.Context, string) (*ticketv1.Event, error)
	approveEventFn     func(context.Context, string) (*ticketv1.Event, error)
}

func (f *fakeTicketClient) CreateOrganizer(ctx context.Context, name, email, phone, slug string) (*ticketv1.Organizer, error) {
	if f.createOrganizerFn != nil {
		return f.createOrganizerFn(ctx, name, email, phone, slug)
	}
	return &ticketv1.Organizer{Id: "org-1", Name: name, Email: email, Phone: phone, Slug: slug, IsActive: true}, nil
}

func (f *fakeTicketClient) GetOrganizerByID(ctx context.Context, id string) (*ticketv1.Organizer, error) {
	if f.getOrganizerByIDFn != nil {
		return f.getOrganizerByIDFn(ctx, id)
	}
	return &ticketv1.Organizer{Id: id, Name: "Organizer One", Email: "organizer@example.com", IsActive: true}, nil
}

func (f *fakeTicketClient) ListOrganizers(ctx context.Context, offset, limit int32) ([]*ticketv1.Organizer, error) {
	if f.listOrganizersFn != nil {
		return f.listOrganizersFn(ctx, offset, limit)
	}
	return []*ticketv1.Organizer{{Id: "org-1", Name: "Organizer One", Email: "organizer@example.com", IsActive: true}}, nil
}

func (f *fakeTicketClient) UpdateOrganizer(ctx context.Context, id, name, phone, slug string) (*ticketv1.Organizer, error) {
	if f.updateOrganizerFn != nil {
		return f.updateOrganizerFn(ctx, id, name, phone, slug)
	}
	return &ticketv1.Organizer{Id: id, Name: name, Phone: phone, Slug: slug, IsActive: true}, nil
}

func (f *fakeTicketClient) CreateEvent(ctx context.Context, organizerID, title, description, venue string, startAt, endAt time.Time, capacity int32) (*ticketv1.Event, error) {
	if f.createEventFn != nil {
		return f.createEventFn(ctx, organizerID, title, description, venue, startAt, endAt, capacity)
	}
	return &ticketv1.Event{Id: "evt-1", OrganizerId: organizerID, Title: title, Venue: venue, Capacity: capacity, Status: ticketv1.EventStatus_EVENT_STATUS_DRAFT}, nil
}

func (f *fakeTicketClient) GetEventByID(ctx context.Context, id string) (*ticketv1.Event, error) {
	if f.getEventByIDFn != nil {
		return f.getEventByIDFn(ctx, id)
	}
	return &ticketv1.Event{Id: id, Title: "Event One", Venue: "HCM", Status: ticketv1.EventStatus_EVENT_STATUS_PUBLISHED}, nil
}

func (f *fakeTicketClient) ListEvents(ctx context.Context, offset, limit int32) ([]*ticketv1.Event, error) {
	if f.listEventsFn != nil {
		return f.listEventsFn(ctx, offset, limit)
	}
	return []*ticketv1.Event{{Id: "evt-1", Title: "Event One", Venue: "HCM", Status: ticketv1.EventStatus_EVENT_STATUS_PUBLISHED}}, nil
}

func (f *fakeTicketClient) PublishEvent(ctx context.Context, id string) (*ticketv1.Event, error) {
	if f.publishEventFn != nil {
		return f.publishEventFn(ctx, id)
	}

	return &ticketv1.Event{Id: id, Title: "Event One", Venue: "HCM", Status: ticketv1.EventStatus_EVENT_STATUS_PUBLISHED}, nil
}

func (f *fakeTicketClient) SubmitEvent(ctx context.Context, id string) (*ticketv1.Event, error) {
	if f.submitEventFn != nil {
		return f.submitEventFn(ctx, id)
	}
	return &ticketv1.Event{Id: id, Status: ticketv1.EventStatus_EVENT_STATUS_PENDING_APPROVAL}, nil
}

func (f *fakeTicketClient) ApproveEvent(ctx context.Context, id string) (*ticketv1.Event, error) {
	if f.approveEventFn != nil {
		return f.approveEventFn(ctx, id)
	}
	return &ticketv1.Event{Id: id, Status: ticketv1.EventStatus_EVENT_STATUS_PUBLISHED}, nil
}

func TestRegisterLoginLogoutFlow(t *testing.T) {
	service := auth.NewTokenService("test-secret")
	identityClient := &fakeIdentityClient{
		createUserFn: func(ctx context.Context, email, password, fullName, phone, role string) (*identityv1.User, error) {
			return &identityv1.User{Id: "u-1", Email: email, FullName: fullName, Phone: phone, IsActive: true}, nil
		},
		authFn: func(ctx context.Context, email, password string) (*identityv1.User, string, error) {
			return &identityv1.User{Id: "u-1", Email: email, FullName: "User One", Phone: "0900", IsActive: true}, "identity-token", nil
		},
		getUserFn: func(ctx context.Context, id string) (*identityv1.User, error) {
			return &identityv1.User{Id: id, Email: "user@example.com", FullName: "User One", IsActive: true}, nil
		},
	}
	ticketClient := &fakeTicketClient{}

	mux := NewMux(service, identityClient, ticketClient)

	registerReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte(`{"email":"user1@example.com","password":"12345678","full_name":"User One","phone":"0900000001"}`)))
	registerRes := httptest.NewRecorder()
	mux.ServeHTTP(registerRes, registerReq)
	if registerRes.Code != http.StatusCreated {
		t.Fatalf("register status = %d; want %d", registerRes.Code, http.StatusCreated)
	}
	if !strings.Contains(registerRes.Body.String(), "user1@example.com") {
		t.Fatalf("register response = %s", registerRes.Body.String())
	}

	loginReq := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"user1@example.com","password":"12345678"}`))
	loginRes := httptest.NewRecorder()
	mux.ServeHTTP(loginRes, loginReq)
	if loginRes.Code != http.StatusOK {
		t.Fatalf("login status = %d; want %d", loginRes.Code, http.StatusOK)
	}

	var loginPayload map[string]any
	if err := json.Unmarshal(loginRes.Body.Bytes(), &loginPayload); err != nil {
		t.Fatalf("login response not json: %v", err)
	}
	accessToken, ok := loginPayload["access_token"].(string)
	if !ok || accessToken == "" {
		t.Fatalf("missing access_token in %+v", loginPayload)
	}

	meReq := httptest.NewRequest(http.MethodGet, "/me", nil)
	meReq.Header.Set("Authorization", "Bearer "+accessToken)
	meRes := httptest.NewRecorder()
	mux.ServeHTTP(meRes, meReq)
	if meRes.Code != http.StatusOK {
		t.Fatalf("me status = %d; want %d", meRes.Code, http.StatusOK)
	}

	logoutReq := httptest.NewRequest(http.MethodPost, "/logout", nil)
	logoutReq.Header.Set("Authorization", "Bearer "+accessToken)
	logoutRes := httptest.NewRecorder()
	mux.ServeHTTP(logoutRes, logoutReq)
	if logoutRes.Code != http.StatusOK {
		t.Fatalf("logout status = %d; want %d", logoutRes.Code, http.StatusOK)
	}
	if !strings.Contains(logoutRes.Body.String(), "logged out") {
		t.Fatalf("logout response = %s", logoutRes.Body.String())
	}
}

func TestTicketOrganizerAndEventRoutes(t *testing.T) {
	service := auth.NewTokenService("test-secret")
	identityClient := &fakeIdentityClient{}
	ticketClient := &fakeTicketClient{
		createOrganizerFn: func(ctx context.Context, name, email, phone, slug string) (*ticketv1.Organizer, error) {
			return &ticketv1.Organizer{Id: "org-1", Name: name, Email: email, Phone: phone, Slug: slug, IsActive: true}, nil
		},
		createEventFn: func(ctx context.Context, organizerID, title, description, venue string, startAt, endAt time.Time, capacity int32) (*ticketv1.Event, error) {
			return &ticketv1.Event{Id: "evt-1", OrganizerId: organizerID, Title: title, Description: description, Venue: venue, Capacity: capacity, Status: ticketv1.EventStatus_EVENT_STATUS_DRAFT}, nil
		},
		listOrganizersFn: func(ctx context.Context, offset, limit int32) ([]*ticketv1.Organizer, error) {
			return []*ticketv1.Organizer{{Id: "org-1", Name: "Zeta", Email: "organizer@example.com", IsActive: true}}, nil
		},
		listEventsFn: func(ctx context.Context, offset, limit int32) ([]*ticketv1.Event, error) {
			return []*ticketv1.Event{{Id: "evt-1", OrganizerId: "org-1", Title: "Launch Night", Venue: "HCM", Status: ticketv1.EventStatus_EVENT_STATUS_DRAFT}}, nil
		},
	}

	mux := NewMux(service, identityClient, ticketClient)
	adminToken, err := service.GenerateWithRole("admin-1", "admin@example.com", "admin")
	if err != nil {
		t.Fatal(err)
	}
	organizerToken, err := service.GenerateWithRole("org-user-1", "organizer@example.com", "organizer")
	if err != nil {
		t.Fatal(err)
	}

	listOrganizersReq := httptest.NewRequest(http.MethodGet, "/organizers?offset=0&limit=20", nil)
	listOrganizersRes := httptest.NewRecorder()
	mux.ServeHTTP(listOrganizersRes, listOrganizersReq)
	if listOrganizersRes.Code != http.StatusOK {
		t.Fatalf("organizers list status = %d; want %d", listOrganizersRes.Code, http.StatusOK)
	}

	createOrganizerReq := httptest.NewRequest(http.MethodPost, "/organizers", strings.NewReader(`{"name":"Zeta","email":"organizer@example.com","phone":"0900000000","slug":"zeta"}`))
	createOrganizerReq.Header.Set("Authorization", "Bearer "+adminToken)
	createOrganizerRes := httptest.NewRecorder()
	mux.ServeHTTP(createOrganizerRes, createOrganizerReq)
	if createOrganizerRes.Code != http.StatusCreated {
		t.Fatalf("organizer create status = %d; want %d", createOrganizerRes.Code, http.StatusCreated)
	}

	listEventsReq := httptest.NewRequest(http.MethodGet, "/events?offset=0&limit=20", nil)
	listEventsRes := httptest.NewRecorder()
	mux.ServeHTTP(listEventsRes, listEventsReq)
	if listEventsRes.Code != http.StatusOK {
		t.Fatalf("events list status = %d; want %d", listEventsRes.Code, http.StatusOK)
	}

	createEventReq := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(`{"organizer_id":"org-1","title":"Launch Night","description":"desc","venue":"HCM","start_at":"2026-10-10T09:00:00Z","end_at":"2026-10-10T11:00:00Z","capacity":100}`))
	createEventReq.Header.Set("Authorization", "Bearer "+organizerToken)
	createEventRes := httptest.NewRecorder()
	mux.ServeHTTP(createEventRes, createEventReq)
	if createEventRes.Code != http.StatusCreated {
		t.Fatalf("event create status = %d; want %d", createEventRes.Code, http.StatusCreated)
	}

	publishReq := httptest.NewRequest(http.MethodPost, "/events/evt-1/publish", nil)
	publishReq.Header.Set("Authorization", "Bearer "+adminToken)
	publishRes := httptest.NewRecorder()
	mux.ServeHTTP(publishRes, publishReq)
	if publishRes.Code != http.StatusOK {
		t.Fatalf("event publish status = %d; want %d", publishRes.Code, http.StatusOK)
	}
}
