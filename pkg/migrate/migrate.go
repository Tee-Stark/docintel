package migrate

import (
	"docintel/pkg/config"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrationConfig holds configuration for database migrations
type MigrationConfig struct {
	MigrationsDir string
	DBConfig      *config.DBConfig
}

// NewMigrationConfig creates a new migration config with default values
func NewMigrationConfig() *MigrationConfig {
	return &MigrationConfig{
		MigrationsDir: "scripts/db/migrations",
		DBConfig:      config.NewDBConfig(),
	}
}

// RunMigrations runs all pending migrations
func (c *MigrationConfig) RunMigrations() error {
	db, err := c.DBConfig.ConnectDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Construct the full path to migrations
	migrationsPath := filepath.Join(cwd, c.MigrationsDir)

	// Verify the directory exists
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory does not exist: %s", migrationsPath)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// CreateNewMigration creates a new migration file
func (c *MigrationConfig) CreateNewMigration(name string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	migrationsPath := filepath.Join(cwd, c.MigrationsDir)

	timestamp := time.Now().UTC().Format("20060102150405")
	upFilename := fmt.Sprintf("%s_%s.up.sql", timestamp, name)
	downFilename := fmt.Sprintf("%s_%s.down.sql", timestamp, name)

	upFile := filepath.Join(migrationsPath, upFilename)
	downFile := filepath.Join(migrationsPath, downFilename)

	// Create the up migration file
	f, err := os.Create(upFile)
	if err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}
	_, err = f.WriteString(fmt.Sprintf("-- Migration: %s (up)\n\n", name))
	if err != nil {
		f.Close()
		return fmt.Errorf("failed to write to up migration file: %w", err)
	}
	f.Close()

	// Create the down migration file
	f, err = os.Create(downFile)
	if err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}
	_, err = f.WriteString(fmt.Sprintf("-- Migration: %s (down)\n\n", name))
	if err != nil {
		f.Close()
		return fmt.Errorf("failed to write to down migration file: %w", err)
	}
	f.Close()

	fmt.Printf("Created new migration files:\n  - %s\n  - %s\n", upFilename, downFilename)
	return nil
}
