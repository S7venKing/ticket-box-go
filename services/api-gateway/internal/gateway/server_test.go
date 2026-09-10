package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	identityv1 "github.com/S7venKing/ticket-box-go/gen/identity/v1"
	"github.com/S7venKing/ticket-box-go/services/api-gateway/internal/auth"
)

type fakeIdentityClient struct {
	createUserFn func(context.Context, string, string, string, string) (*identityv1.User, error)
	authFn       func(context.Context, string, string) (*identityv1.User, string, error)
	getUserFn    func(context.Context, string) (*identityv1.User, error)
}

func (f *fakeIdentityClient) CreateUser(ctx context.Context, email, password, fullName, phone string) (*identityv1.User, error) {
	if f.createUserFn != nil {
		return f.createUserFn(ctx, email, password, fullName, phone)
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

func TestRegisterLoginLogoutFlow(t *testing.T) {
	service := auth.NewTokenService("test-secret")
	identityClient := &fakeIdentityClient{
		createUserFn: func(ctx context.Context, email, password, fullName, phone string) (*identityv1.User, error) {
			return &identityv1.User{Id: "u-1", Email: email, FullName: fullName, Phone: phone, IsActive: true}, nil
		},
		authFn: func(ctx context.Context, email, password string) (*identityv1.User, string, error) {
			return &identityv1.User{Id: "u-1", Email: email, FullName: "User One", Phone: "0900", IsActive: true}, "identity-token", nil
		},
		getUserFn: func(ctx context.Context, id string) (*identityv1.User, error) {
			return &identityv1.User{Id: id, Email: "user@example.com", FullName: "User One", IsActive: true}, nil
		},
	}

	mux := NewMux(service, identityClient)

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
