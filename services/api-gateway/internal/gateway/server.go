package gateway

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/S7venKing/ticket-box-go/services/api-gateway/internal/auth"
)

func RunHTTPServer(tokenService *auth.TokenService, identityClient *IdentityClient) {
	mux := http.NewServeMux()

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

		user, err := identityClient.CreateUser(context.Background(), req.Email, req.Password, req.FullName, req.Phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]any{
			"user": user,
		}); err != nil {
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

		jwt, err := tokenService.Generate(user.GetId(), user.GetEmail())
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

	mux.HandleFunc("/me", AuthMiddleware(tokenService, protected))
	mux.HandleFunc("/profile", AuthMiddleware(tokenService, protected))
	mux.HandleFunc("/users/me", AuthMiddleware(tokenService, protected))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("api-gateway listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
