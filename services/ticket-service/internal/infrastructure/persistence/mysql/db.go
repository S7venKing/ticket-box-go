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

func NewDB(ctx context.Context, cfg Config) (*sql.DB, error) {
	if err := ensureDatabaseSchema(ctx, cfg); err != nil {
		return nil, err
	}

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
	driverCfg.Collation = "utf8mb4_general_ci"

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

func ensureDatabaseSchema(ctx context.Context, cfg Config) error {
	rootCfg := mysqldriver.NewConfig()
	rootCfg.User = cfg.User
	rootCfg.Passwd = cfg.Password
	rootCfg.Net = "tcp"
	rootCfg.Addr = net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	rootCfg.ParseTime = true
	rootCfg.Loc = time.UTC
	rootCfg.Timeout = 5 * time.Second
	rootCfg.ReadTimeout = 30 * time.Second
	rootCfg.WriteTimeout = 30 * time.Second
	rootCfg.Collation = "utf8mb4_general_ci"

	connector, err := mysqldriver.NewConnector(rootCfg)
	if err != nil {
		return fmt.Errorf("mysql root connector: %w", err)
	}

	db := sql.OpenDB(connector)
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping mysql: %w", err)
	}

	createDatabaseSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Database)
	if _, err := db.ExecContext(ctx, createDatabaseSQL); err != nil {
		return fmt.Errorf("create database %s: %w", cfg.Database, err)
	}

	organizerTableSQL := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS `%s`.`organizers` ("+
			"id CHAR(36) NOT NULL,"+
			"email VARCHAR(255) NOT NULL,"+
			"name VARCHAR(255) NOT NULL,"+
			"phone VARCHAR(50) NULL,"+
			"slug VARCHAR(255) NOT NULL,"+
			"is_active BOOLEAN NOT NULL DEFAULT TRUE,"+
			"created_at DATETIME(6) NOT NULL,"+
			"updated_at DATETIME(6) NOT NULL,"+
			"PRIMARY KEY (id),"+
			"UNIQUE KEY ux_organizers_email (email),"+
			"UNIQUE KEY ux_organizers_slug (slug)"+
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;",
		cfg.Database,
	)
	if _, err := db.ExecContext(ctx, organizerTableSQL); err != nil {
		return fmt.Errorf("create organizers table: %w", err)
	}

	eventTableSQL := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS `%s`.`events` ("+
			"id CHAR(36) NOT NULL,"+
			"organizer_id CHAR(36) NOT NULL,"+
			"title VARCHAR(255) NOT NULL,"+
			"description TEXT NULL,"+
			"venue VARCHAR(255) NOT NULL,"+
			"start_at DATETIME(6) NOT NULL,"+
			"end_at DATETIME(6) NOT NULL,"+
			"capacity INT NOT NULL,"+
			"status VARCHAR(32) NOT NULL DEFAULT 'draft',"+
			"created_at DATETIME(6) NOT NULL,"+
			"updated_at DATETIME(6) NOT NULL,"+
			"PRIMARY KEY (id),"+
			"KEY ix_events_organizer_created (organizer_id, created_at),"+
			"KEY ix_events_status_created (status, created_at),"+
			"CONSTRAINT fk_events_organizer FOREIGN KEY (organizer_id) REFERENCES `%s`.`organizers`(id)"+
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;",
		cfg.Database,
		cfg.Database,
	)
	if _, err := db.ExecContext(ctx, eventTableSQL); err != nil {
		return fmt.Errorf("create events table: %w", err)
	}

	return nil
}
