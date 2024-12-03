package pg

import (
	"currency_eval/internal/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	logger *zap.Logger
	conn   *sqlx.DB
	cfg    config.Config
}

func PostgresDSN(c config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDB,
	)
}

func NewPostgresRepository(logger *zap.Logger, config config.Config) (*PostgresRepository, error) {
	dsn := PostgresDSN(config)
	logger.Info("Postgres DSN", zap.String("dsn", dsn))
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		logger.Error("Failed to connect to PostgreSQL", zap.Error(err))
		return nil, err
	}

	if err := Migrate(db.DB, config, logger); err != nil {
		logger.Error("Failed to run migrations", zap.Error(err))
		return nil, err
	}

	repo := PostgresRepository{
		logger: logger,
		conn:   db,
		cfg:    config,
	}
	return &repo, nil
}
