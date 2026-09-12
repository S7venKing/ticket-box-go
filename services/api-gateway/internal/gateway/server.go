package gateway

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	identityv1 "github.com/S7venKing/ticket-box-go/gen/identity/v1"
	ticketv1 "github.com/S7venKing/ticket-box-go/gen/ticket/v1"
	"github.com/S7venKing/ticket-box-go/services/api-gateway/internal/auth"
)

type IdentityGatewayClient interface {
	CreateUser(ctx context.Context, email, password, fullName, phone, role string) (*identityv1.User, error)
	AuthenticateUser(ctx context.Context, email, password string) (*identityv1.User, string, error)
	GetUserByID(ctx context.Context, id string) (*identityv1.User, error)
}

type IdentityUserListClient interface {
	ListUsers(ctx context.Context, offset, limit int32) ([]*identityv1.User, error)
}
type IdentityMembershipClient interface {
	AssignUserToOrganizer(ctx context.Context, userID, organizerID string) (*identityv1.User, error)
}
type IdentityOrganizerClient interface {
	CreateOrganizer(ctx context.Context, name, email, phone, slug string) (*identityv1.Organizer, error)
	GetOrganizerByID(ctx context.Context, id string) (*identityv1.Organizer, error)
	ListOrganizers(ctx context.Context, offset, limit int32) ([]*identityv1.Organizer, error)
	UpdateOrganizer(ctx context.Context, id, name, email, phone, slug string) (*identityv1.Organizer, error)
}
type TicketGatewayClient interface {
	CreateEvent(ctx context.Context, organizerID, title, description, venue string, startAt, endAt time.Time, capacity int32) (*ticketv1.Event, error)
	GetEventByID(ctx context.Context, id string) (*ticketv1.Event, error)
	ListEvents(ctx context.Context, offset, limit int32) ([]*ticketv1.Event, error)
	PublishEvent(ctx context.Context, id string) (*ticketv1.Event, error)
	SubmitEvent(ctx context.Context, id string) (*ticketv1.Event, error)
	ApproveEvent(ctx context.Context, id string) (*ticketv1.Event, error)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Requested-With")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewMux(tokenService *auth.TokenService, identityClient IdentityGatewayClient, ticketClient TicketGatewayClient) *http.ServeMux {
	mux := http.NewServeMux()
	organizerClient, organizerAvailable := identityClient.(IdentityOrganizerClient)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			FullName string `json:"full_name"`
			Phone    string `json:"phone"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, err := identityClient.CreateUser(context.Background(), req.Email, req.Password, req.FullName, req.Phone, "user")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]any{"user": user}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, tokenFromIdentity, err := identityClient.AuthenticateUser(context.Background(), req.Email, req.Password)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		jwt, err := tokenService.GenerateWithRole(user.GetId(), user.GetEmail(), user.GetRole())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"user":           user,
			"access_token":   jwt,
			"identity_token": tokenFromIdentity,
		}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}
			if _, err := tokenService.Validate(strings.TrimSpace(parts[1])); err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"message": "logged out"}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})

	protected := func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDContextKey{}).(string)
		userEmail, _ := r.Context().Value(userEmailContextKey{}).(string)
		if userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := identityClient.GetUserByID(r.Context(), userID)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"user_id":    userID,
			"user_email": userEmail,
			"user":       user,
		}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}

	requireAdmin := func(next http.HandlerFunc) http.HandlerFunc {
		return AuthMiddleware(tokenService, func(w http.ResponseWriter, r *http.Request) {
			role, _ := r.Context().Value(userRoleContextKey{}).(string)
			if role != "admin" {
				http.Error(w, "admin role required", http.StatusForbidden)
				return
			}
			next(w, r)
		})
	}

	organizerOwned := func(r *http.Request, organizerID string) bool {
		role, _ := r.Context().Value(userRoleContextKey{}).(string)
		if role == "admin" {
			return true
		}
		userID, _ := r.Context().Value(userIDContextKey{}).(string)
		user, err := identityClient.GetUserByID(r.Context(), userID)
		return err == nil && user.GetOrganizerId() != "" && user.GetOrganizerId() == organizerID
	}

	requireOrganizer := func(next http.HandlerFunc) http.HandlerFunc {
		return AuthMiddleware(tokenService, func(w http.ResponseWriter, r *http.Request) {
			role, _ := r.Context().Value(userRoleContextKey{}).(string)
			if role == "admin" {
				next(w, r)
				return
			}
			userID, _ := r.Context().Value(userIDContextKey{}).(string)
			user, err := identityClient.GetUserByID(r.Context(), userID)
			if err != nil || user.GetOrganizerId() == "" {
				http.Error(w, "organizer membership required", http.StatusForbidden)
				return
			}
			next(w, r)
		})
	}

	mux.HandleFunc("/me", AuthMiddleware(tokenService, protected))
	mux.HandleFunc("/profile", AuthMiddleware(tokenService, protected))
	mux.HandleFunc("/users/me", AuthMiddleware(tokenService, protected))
	mux.HandleFunc("GET /admin/users", requireAdmin(func(w http.ResponseWriter, r *http.Request) {
		listClient, ok := identityClient.(IdentityUserListClient)
		if !ok {
			http.Error(w, "user listing is unavailable", http.StatusInternalServerError)
			return
		}
		users, err := listClient.ListUsers(r.Context(), 0, 100)
		if err != nil {
			http.Error(w, "failed to load users", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"users": users}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}))

	mux.HandleFunc("POST /admin/organizers", requireAdmin(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Phone    string `json:"phone"`
			Slug     string `json:"slug"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		account, err := identityClient.CreateUser(r.Context(), req.Email, req.Password, req.Name, req.Phone, "organizer")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !organizerAvailable {
			http.Error(w, "organizer management unavailable", http.StatusNotImplemented)
			return
		}
		organizer, err := organizerClient.CreateOrganizer(r.Context(), req.Name, req.Email, req.Phone, req.Slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if memberClient, ok := identityClient.(IdentityMembershipClient); ok {
			if _, err := memberClient.AssignUserToOrganizer(r.Context(), account.GetId(), organizer.GetId()); err != nil {
				http.Error(w, "failed to assign organizer membership", http.StatusBadGateway)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"account": account, "organizer": organizer})
	}))
	mux.HandleFunc("POST /admin/organizers/{id}/members", requireAdmin(func(w http.ResponseWriter, r *http.Request) {
		memberClient, ok := identityClient.(IdentityMembershipClient)
		if !ok {
			http.Error(w, "membership management unavailable", http.StatusNotImplemented)
			return
		}
		var req struct {
			UserID string `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		user, err := memberClient.AssignUserToOrganizer(r.Context(), req.UserID, r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"user": user})
	}))

	mux.HandleFunc("GET /organizers", func(w http.ResponseWriter, r *http.Request) {
		if !organizerAvailable {
			http.Error(w, "organizer management unavailable", http.StatusNotImplemented)
			return
		}
		offset, limit := parseOffsetLimit(r)
		organizers, err := organizerClient.ListOrganizers(r.Context(), int32(offset), int32(limit))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"organizers": organizers}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("POST /organizers", requireAdmin(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Phone string `json:"phone"`
			Slug  string `json:"slug"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if !organizerAvailable {
			http.Error(w, "organizer management unavailable", http.StatusNotImplemented)
			return
		}
		organizer, err := organizerClient.CreateOrganizer(r.Context(), req.Name, req.Email, req.Phone, req.Slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]any{"organizer": organizer}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}))
	mux.HandleFunc("GET /organizers/{id}", func(w http.ResponseWriter, r *http.Request) {
		if !organizerAvailable {
			http.Error(w, "organizer management unavailable", http.StatusNotImplemented)
			return
		}
		id := r.PathValue("id")
		organizer, err := organizerClient.GetOrganizerByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"organizer": organizer}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("PUT /organizers/{id}", requireAdmin(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var req struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Phone string `json:"phone"`
			Slug  string `json:"slug"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if !organizerAvailable {
			http.Error(w, "organizer management unavailable", http.StatusNotImplemented)
			return
		}
		organizer, err := organizerClient.UpdateOrganizer(r.Context(), id, req.Name, req.Email, req.Phone, req.Slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"organizer": organizer}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}))

	listEvents := func(w http.ResponseWriter, r *http.Request) {
		offset, limit := parseOffsetLimit(r)
		events, err := ticketClient.ListEvents(r.Context(), int32(offset), int32(limit))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		role, _ := r.Context().Value(userRoleContextKey{}).(string)
		if role != "" && role != "admin" {
			userID, _ := r.Context().Value(userIDContextKey{}).(string)
			user, userErr := identityClient.GetUserByID(r.Context(), userID)
			if userErr != nil || user.GetOrganizerId() == "" {
				http.Error(w, "organizer membership required", http.StatusForbidden)
				return
			}
			filtered := events[:0]
			for _, event := range events {
				if event.GetOrganizerId() == user.GetOrganizerId() {
					filtered = append(filtered, event)
				}
			}
			events = filtered
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"events": events}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
	mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
		if strings.TrimSpace(r.Header.Get("Authorization")) == "" {
			listEvents(w, r)
			return
		}
		AuthMiddleware(tokenService, listEvents)(w, r)
	})
	mux.HandleFunc("POST /events", requireOrganizer(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			OrganizerID string `json:"organizer_id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Venue       string `json:"venue"`
			StartAt     string `json:"start_at"`
			EndAt       string `json:"end_at"`
			Capacity    int32  `json:"capacity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if !organizerOwned(r, req.OrganizerID) {
			http.Error(w, "organizer does not own this account", http.StatusForbidden)
			return
		}
		startAt, err := parseTime(req.StartAt)
		if err != nil {
			http.Error(w, "invalid start_at", http.StatusBadRequest)
			return
		}
		endAt, err := parseTime(req.EndAt)
		if err != nil {
			http.Error(w, "invalid end_at", http.StatusBadRequest)
			return
		}
		event, err := ticketClient.CreateEvent(r.Context(), req.OrganizerID, req.Title, req.Description, req.Venue, startAt, endAt, req.Capacity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]any{"event": event}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}))
	mux.HandleFunc("GET /events/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		event, err := ticketClient.GetEventByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if strings.TrimSpace(r.Header.Get("Authorization")) != "" {
			claims, tokenErr := tokenService.Validate(strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer")))
			if tokenErr != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			if claims.Role != "admin" {
				user, userErr := identityClient.GetUserByID(r.Context(), claims.UserID)
				if userErr != nil || user.GetOrganizerId() == "" || user.GetOrganizerId() != event.GetOrganizerId() {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
				}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"event": event}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("POST /events/{id}/submit", requireOrganizer(func(w http.ResponseWriter, r *http.Request) {
		item, err := ticketClient.GetEventByID(r.Context(), r.PathValue("id"))
		if err != nil || !organizerOwned(r, item.GetOrganizerId()) {
			http.Error(w, "organizer does not own this event", http.StatusForbidden)
			return
		}
		event, err := ticketClient.SubmitEvent(r.Context(), r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"event": event})
	}))
	mux.HandleFunc("POST /events/{id}/approve", requireAdmin(func(w http.ResponseWriter, r *http.Request) {
		event, err := ticketClient.ApproveEvent(r.Context(), r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"event": event})
	}))
	mux.HandleFunc("POST /events/{id}/publish", requireAdmin(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		event, err := ticketClient.PublishEvent(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"event": event}); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}))

	return mux
}

func parseOffsetLimit(r *http.Request) (int, int) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20
	}
	return offset, limit
}

func parseTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, value)
}

func RunHTTPServer(tokenService *auth.TokenService, identityClient *IdentityClient, ticketClient *TicketClient) {
	mux := NewMux(tokenService, identityClient, ticketClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("api-gateway listening on :%s", port)
	if err := http.ListenAndServe(":"+port, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}
