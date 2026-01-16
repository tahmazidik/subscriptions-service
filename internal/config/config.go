package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort string

	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string
}

func Load() Config {
	cfg := Config{
		AppPort:   getenv("APP_PORT", "8080"),
		DBHost:    getenv("DB_HOST", "localhost"),
		DBPort:    getenv("DB_PORT", "5432"),
		DBUser:    getenv("DB_USER", "app"),
		DBPass:    getenv("DB_PASSWORD", "app"),
		DBName:    getenv("DB_NAME", "app"),
		DBSSLMode: getenv("DB_SSLMODE", "disable"),
	}
	return cfg
}

func (c Config) DSN() string {
	// Формат для pgx: postgres://user:pass@host:port/dbname?sslmode=disable
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
