package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	ticketv1 "github.com/S7venKing/ticket-box-go/gen/ticket/v1"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/command"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/application/query"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/config"
	"github.com/S7venKing/ticket-box-go/services/ticket-service/internal/infrastructure/persistence/mysql"
	ticketgrpc "github.com/S7venKing/ticket-box-go/services/ticket-service/internal/interfaces/grpc"
)

func main() {
	if err := run(); err != nil {
		slog.Error("ticket-service exited with error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		slog.String("service", cfg.AppName),
		slog.String("env", cfg.AppEnv),
	)
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := mysql.NewDB(ctx, mysql.Config{
		Host:     cfg.MySQLHost,
		Port:     cfg.MySQLPort,
		Database: cfg.MySQLDatabase,
		User:     cfg.MySQLUser,
		Password: cfg.MySQLPassword,
	})
	if err != nil {
		return err
	}
	defer db.Close()

	organizerRepository := mysql.NewOrganizerRepository(db)
	eventRepository := mysql.NewEventRepository(db)

	ticketServer := ticketgrpc.NewTicketServer(
		logger,
		command.NewCreateOrganizerHandler(organizerRepository),
		query.NewGetOrganizerByIDHandler(organizerRepository),
		query.NewListOrganizersHandler(organizerRepository),
		command.NewUpdateOrganizerHandler(organizerRepository),
		command.NewCreateEventHandler(eventRepository),
		query.NewGetEventByIDHandler(eventRepository),
		query.NewListEventsHandler(eventRepository),
		command.NewPublishEventHandler(eventRepository),
		command.NewSubmitEventHandler(eventRepository),
		command.NewApproveEventHandler(eventRepository),
	)

	grpcServer := grpc.NewServer(
		grpc.ConnectionTimeout(10*time.Second),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              30 * time.Second,
			Timeout:           10 * time.Second,
		}),
	)

	ticketv1.RegisterTicketServiceServer(grpcServer, ticketServer)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(ticketv1.TicketService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	if !cfg.IsProduction() {
		reflection.Register(grpcServer)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- grpcServer.Serve(listener)
	}()

	logger.Info("ticket gRPC server started", slog.String("addr", listener.Addr().String()))

	select {
	case err := <-serveErr:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}

	logger.Info("shutting down", slog.Duration("timeout", cfg.ShutdownTimeout))
	healthServer.Shutdown()

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-time.After(cfg.ShutdownTimeout):
		logger.Warn("graceful stop timed out, forcing stop")
		grpcServer.Stop()
	}

	logger.Info("stopped")
	return nil
}
