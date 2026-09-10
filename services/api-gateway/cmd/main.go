package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/S7venKing/ticket-box-go/services/api-gateway/internal/auth"
	"github.com/S7venKing/ticket-box-go/services/api-gateway/internal/gateway"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "change-me-in-production"
	}

	identityAddr := os.Getenv("IDENTITY_SERVICE_ADDR")
	if identityAddr == "" {
		identityAddr = "localhost:50051"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, identityAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial identity service: %v", err)
	}
	defer conn.Close()

	identityClient := gateway.NewIdentityClient(conn)
	tokenService := auth.NewTokenService(jwtSecret)

	gateway.RunHTTPServer(tokenService, identityClient)
}
