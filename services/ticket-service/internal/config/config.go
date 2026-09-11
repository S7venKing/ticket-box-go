package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppName         string
	AppEnv          string
	GRPCPort        int
	ShutdownTimeout time.Duration
	JWTSecret       string
	JWTTTL          time.Duration

	MySQLHost     string
	MySQLPort     int
	MySQLDatabase string
	MySQLUser     string
	MySQLPassword string
}

func (c Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func Load() (Config, error) {
	cfg := Config{
		AppName:       getEnv("APP_NAME", "ticket-service"),
		AppEnv:        getEnv("APP_ENV", "development"),
		JWTSecret:     getEnv("JWT_SECRET", "s7venKing@123"),
		MySQLHost:     getEnv("MYSQL_HOST", "localhost"),
		MySQLDatabase: getEnv("MYSQL_DATABASE", "ticket_db"),
		MySQLUser:     getEnv("MYSQL_USER", "root"),
		MySQLPassword: os.Getenv("MYSQL_PASSWORD"),
	}

	var err error
	if cfg.GRPCPort, err = getEnvInt("APP_PORT", 50052); err != nil {
		return Config{}, err
	}
	if cfg.MySQLPort, err = getEnvInt("MYSQL_PORT", 3306); err != nil {
		return Config{}, err
	}
	if cfg.ShutdownTimeout, err = getEnvDuration("SHUTDOWN_TIMEOUT", 15*time.Second); err != nil {
		return Config{}, err
	}
	if cfg.JWTTTL, err = getEnvDuration("JWT_TTL", 30*time.Minute); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("config: %s must be an integer: %w", key, err)
	}
	return n, nil
}

func getEnvDuration(key string, fallback time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("config: %s must be a duration like 15s: %w", key, err)
	}
	return d, nil
}
