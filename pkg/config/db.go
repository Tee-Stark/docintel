package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	// DatabaseURL overrides individual fields when set (e.g. Fly.io DATABASE_URL).
	DatabaseURL string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Host:        getEnvOrDefault("DB_HOST", "localhost"),
		Port:        getEnvOrDefault("DB_PORT", "5432"),
		User:        getEnvOrDefault("DB_USER", "postgres"),
		Password:    getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:      getEnvOrDefault("DB_NAME", "docintel"),
	}
}

func (c *DBConfig) ConnectDB() (*sql.DB, error) {
	var dsn string
	if c.DatabaseURL != "" {
		dsn = c.DatabaseURL
		if len(dsn) > 0 && !containsParam(dsn, "sslmode") {
			if containsParam(dsn, "?") {
				dsn += "&sslmode=disable"
			} else {
				dsn += "?sslmode=disable"
			}
		}
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Password, c.DBName)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(3 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// containsParam reports whether s contains the substring sub, used to check
// whether a DSN or URL already has a given parameter.
func containsParam(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func (c *DBConfig) ConnectPgxPool() (*pgxpool.Pool, error) {
	var dsn string
	if c.DatabaseURL != "" {
		dsn = c.DatabaseURL
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Password, c.DBName)
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
