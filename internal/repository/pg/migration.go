package pg

import (
	"currency_eval/internal/config"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
	"path/filepath"
)

func Migrate(db *sql.DB, cfg config.Config, logger *zap.Logger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to get pg driver %v", err)
	}
	absPath, _ := filepath.Abs(".")

	logger.Debug("Current working directory", zap.String("path", absPath))
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+absPath+"/internal/repository/pg/migrations",
		cfg.PostgresDB,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to get migrate instance %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate %v", err)
	}

	return nil
}
