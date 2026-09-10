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

		user, token, err := identityClient.AuthenticateUser(context.Background(), req.Email, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		jwt, err := tokenService.Generate(user.GetId(), user.GetEmail())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_ = token
		_ = json.NewEncoder(w).Encode(map[string]any{
			"user":         user,
			"access_token": jwt,
		})
	})

	mux.HandleFunc("/me", AuthMiddleware(tokenService, func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("user_id").(string)
		userEmail, _ := r.Context().Value("user_email").(string)
		_ = userID
		_ = userEmail
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"user_id": userID, "email": userEmail})
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("api-gateway listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
