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

	ticketAddr := os.Getenv("TICKET_SERVICE_ADDR")
	if ticketAddr == "" {
		ticketAddr = "localhost:50052"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	identityConn, err := grpc.DialContext(ctx, identityAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial identity service: %v", err)
	}
	defer identityConn.Close()

	ticketConn, err := grpc.DialContext(ctx, ticketAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial ticket service: %v", err)
	}
	defer ticketConn.Close()

	identityClient := gateway.NewIdentityClient(identityConn)
	ticketClient := gateway.NewTicketClient(ticketConn)
	tokenService := auth.NewTokenService(jwtSecret)

	gateway.RunHTTPServer(tokenService, identityClient, ticketClient)
}
