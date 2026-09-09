package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

// NewDB opens a pooled connection and verifies it with a bounded ping.
func NewDB(ctx context.Context, cfg Config) (*sql.DB, error) {
	driverCfg := mysqldriver.NewConfig()
	driverCfg.User = cfg.User
	driverCfg.Passwd = cfg.Password
	driverCfg.Net = "tcp"
	driverCfg.Addr = net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	driverCfg.DBName = cfg.Database
	driverCfg.ParseTime = true
	driverCfg.Loc = time.UTC
	driverCfg.Timeout = 5 * time.Second
	driverCfg.ReadTimeout = 30 * time.Second
	driverCfg.WriteTimeout = 30 * time.Second
	driverCfg.Params = map[string]string{"charset": "utf8mb4"}

	connector, err := mysqldriver.NewConnector(driverCfg)
	if err != nil {
		return nil, fmt.Errorf("mysql connector: %w", err)
	}

	db := sql.OpenDB(connector)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	return db, nil
}
