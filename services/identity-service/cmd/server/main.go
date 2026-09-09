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

	identityv1 "github.com/S7venKing/ticket-box-go/gen/identity/v1"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/command"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/query"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/config"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/infrastructure/persistence/mysql"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/infrastructure/security/password"
	identitygrpc "github.com/S7venKing/ticket-box-go/services/identity-service/internal/interfaces/grpc"
)

func main() {
	if err := run(); err != nil {
		slog.Error("identity-service exited with error", slog.Any("error", err))
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

	// ---- infrastructure ----
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

	userRepository := mysql.NewUserRepository(db)
	passwordHasher := password.NewArgon2Hasher()

	// ---- application + inbound adapter ----
	identityServer := identitygrpc.NewIdentityServer(
		logger,
		command.NewCreateUserHandler(userRepository, passwordHasher),
		query.NewGetUserByIDHandler(userRepository),
		query.NewGetUserByEmailHandler(userRepository),
		command.NewUpdateProfileHandler(userRepository),
		command.NewActivateUserHandler(userRepository),
		command.NewDeactivateUserHandler(userRepository),
	)

	// ---- gRPC server ----
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(identitygrpc.UnaryInterceptors(logger)...),
		grpc.ConnectionTimeout(10*time.Second),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              30 * time.Second,
			Timeout:           10 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	)

	identityv1.RegisterIdentityServiceServer(grpcServer, identityServer)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(identityv1.IdentityService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)

	if !cfg.IsProduction() {
		reflection.Register(grpcServer) // grpcurl / evans in dev only
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- grpcServer.Serve(listener)
	}()

	logger.Info("grpc server started", slog.String("addr", listener.Addr().String()))

	select {
	case err := <-serveErr:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}

	// ---- graceful shutdown ----
	logger.Info("shutting down", slog.Duration("timeout", cfg.ShutdownTimeout))

	healthServer.Shutdown() // NOT_SERVING first, so load balancers stop routing to us

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
