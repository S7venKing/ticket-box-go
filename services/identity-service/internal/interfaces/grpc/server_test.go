package grpc_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	identityv1 "github.com/S7venKing/ticket-box-go/gen/identity/v1"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/command"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/application/query"
	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
	identitygrpc "github.com/S7venKing/ticket-box-go/services/identity-service/internal/interfaces/grpc"
)

// memRepo is an in-memory user.Repository so the gRPC adapter can be exercised
// end-to-end (interceptors -> server -> handlers -> error mapping) without MySQL.
type memRepo struct {
	mu    sync.Mutex
	users map[uuid.UUID]*user.User
}

func newMemRepo() *memRepo { return &memRepo{users: map[uuid.UUID]*user.User{}} }

func (r *memRepo) Create(_ context.Context, u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.users {
		if existing.Email == u.Email {
			return user.ErrEmailAlreadyExists
		}
	}

	c := *u
	r.users[u.ID] = &c

	return nil
}

func (r *memRepo) GetByID(_ context.Context, id uuid.UUID) (*user.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[id]
	if !ok {
		return nil, user.ErrUserNotFound
	}

	c := *u

	return &c, nil
}

func (r *memRepo) GetByEmail(_ context.Context, email string) (*user.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Email == email {
			c := *u
			return &c, nil
		}
	}

	return nil, user.ErrUserNotFound
}

func (r *memRepo) Update(_ context.Context, u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[u.ID]; !ok {
		return user.ErrUserNotFound
	}

	c := *u
	r.users[u.ID] = &c

	return nil
}

func (r *memRepo) ExistsByEmail(_ context.Context, email string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Email == email {
			return true, nil
		}
	}

	return false, nil
}

type plainHasher struct{}

func (plainHasher) Hash(p string) (string, error) { return "hash:" + p, nil }

func (plainHasher) Compare(p, h string) error {
	if "hash:"+p != h {
		return errors.New("mismatch")
	}

	return nil
}

func newTestClient(t *testing.T) identityv1.IdentityServiceClient {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := newMemRepo()

	server := identitygrpc.NewIdentityServer(
		logger,
		"super-secret",
		command.NewCreateUserHandler(repo, plainHasher{}),
		command.NewAuthenticateUserHandler(repo, plainHasher{}),
		query.NewGetUserByIDHandler(repo),
		query.NewGetUserByEmailHandler(repo),
		command.NewUpdateProfileHandler(repo),
		command.NewActivateUserHandler(repo),
		command.NewDeactivateUserHandler(repo),
	)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(identitygrpc.UnaryInterceptors(logger)...))
	identityv1.RegisterIdentityServiceServer(grpcServer, server)

	lis := bufconn.Listen(1 << 20)
	go func() { _ = grpcServer.Serve(lis) }()
	t.Cleanup(grpcServer.Stop)

	conn, err := grpc.NewClient(
		"passthrough:///bufconn",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	return identityv1.NewIdentityServiceClient(conn)
}

func TestIdentityServer_EndToEnd(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := newTestClient(t)

	created, err := client.CreateUser(ctx, &identityv1.CreateUserRequest{
		Email: " Alice@Example.com ", Password: "correct-horse", FullName: "Alice", Phone: "0900000000",
	})
	require.NoError(t, err)
	require.Equal(t, "alice@example.com", created.GetUser().GetEmail())
	require.True(t, created.GetUser().GetIsActive())
	require.WithinDuration(t, time.Now(), created.GetUser().GetCreatedAt().AsTime(), time.Minute)

	// duplicate email -> AlreadyExists
	_, err = client.CreateUser(ctx, &identityv1.CreateUserRequest{Email: "alice@example.com", Password: "correct-horse", FullName: "Alice"})
	require.Equal(t, codes.AlreadyExists, status.Code(err))

	// short password -> InvalidArgument
	_, err = client.CreateUser(ctx, &identityv1.CreateUserRequest{Email: "bob@example.com", Password: "short", FullName: "Bob"})
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	// malformed uuid -> InvalidArgument
	_, err = client.GetUserById(ctx, &identityv1.GetUserByIdRequest{Id: "not-a-uuid"})
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	// unknown user -> NotFound
	_, err = client.GetUserById(ctx, &identityv1.GetUserByIdRequest{Id: uuid.NewString()})
	require.Equal(t, codes.NotFound, status.Code(err))

	// authentication success returns a JWT token
	authRes, err := client.AuthenticateUser(ctx, &identityv1.AuthenticateUserRequest{
		Email:    "alice@example.com",
		Password: "correct-horse",
	})
	require.NoError(t, err)
	require.NotEmpty(t, authRes.GetAccessToken())
	require.Equal(t, "alice@example.com", authRes.GetUser().GetEmail())

	// wrong password -> Unauthenticated
	_, err = client.AuthenticateUser(ctx, &identityv1.AuthenticateUserRequest{
		Email:    "alice@example.com",
		Password: "wrong-password",
	})
	require.Equal(t, codes.Unauthenticated, status.Code(err))

	// deactivate twice -> FailedPrecondition
	id := created.GetUser().GetId()
	_, err = client.DeactivateUser(ctx, &identityv1.DeactivateUserRequest{Id: id})
	require.NoError(t, err)
	_, err = client.DeactivateUser(ctx, &identityv1.DeactivateUserRequest{Id: id})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))

	// lookup by email is case-insensitive
	byEmail, err := client.GetUserByEmail(ctx, &identityv1.GetUserByEmailRequest{Email: "ALICE@example.com"})
	require.NoError(t, err)
	require.Equal(t, id, byEmail.GetUser().GetId())
	require.False(t, byEmail.GetUser().GetIsActive())
}
