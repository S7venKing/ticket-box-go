//go:build integration

// Integration tests need a MySQL with the users table (see migrations/).
//
//	docker compose up -d mysql
//	go test -tags integration ./test/integration/...
package integration

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/S7venKing/ticket-box-go/services/identity-service/internal/domain/user"
	mysqlrepo "github.com/S7venKing/ticket-box-go/services/identity-service/internal/infrastructure/persistence/mysql"
)

func TestUserRepository_CreateGetDuplicate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	port, _ := strconv.Atoi(getenv("MYSQL_PORT", "3306"))

	db, err := mysqlrepo.NewDB(ctx, mysqlrepo.Config{
		Host:     getenv("MYSQL_HOST", "localhost"),
		Port:     port,
		Database: getenv("MYSQL_DATABASE", "identity_db"),
		User:     getenv("MYSQL_USER", "root"),
		Password: getenv("MYSQL_PASSWORD", "root"),
	})
	require.NoError(t, err)
	defer db.Close()

	repo := mysqlrepo.NewUserRepository(db)

	email := fmt.Sprintf("it-%s@example.com", uuid.NewString())

	u, err := user.NewUser(email, "hash", "Integration Test", "")
	require.NoError(t, err)
	require.NoError(t, repo.Create(ctx, u))

	got, err := repo.GetByID(ctx, u.ID)
	require.NoError(t, err)
	require.Equal(t, email, got.Email)

	exists, err := repo.ExistsByEmail(ctx, email)
	require.NoError(t, err)
	require.True(t, exists)

	dup, err := user.NewUser(email, "hash", "Duplicate", "")
	require.NoError(t, err)
	require.ErrorIs(t, repo.Create(ctx, dup), user.ErrEmailAlreadyExists)

	_, err = repo.GetByID(ctx, uuid.New())
	require.ErrorIs(t, err, user.ErrUserNotFound)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
