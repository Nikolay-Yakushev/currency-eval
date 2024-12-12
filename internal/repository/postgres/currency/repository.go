package currency

import (
	"context"
	"currency_eval/internal/models"
	"currency_eval/internal/repository/postgres"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type PostgresCurrencyRepository struct {
	logger *zap.Logger
	conn   *sqlx.DB
	cfg    postgres.Config
}

func NewCurrencyRepository(logger *zap.Logger, config postgres.Config) (*PostgresCurrencyRepository, error) {
	var dsn = config.DSN()

	logger.Debug("Postgres DSN", zap.String("dsn", dsn))

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres currency repository, %w", err)
	}

	if err := postgres.Migrate(db.DB, config, logger); err != nil {
		return nil, fmt.Errorf("failed to run migrations, %w", err)
	}

	repo := PostgresCurrencyRepository{
		logger: logger,
		conn:   db,
		cfg:    config,
	}
	return &repo, nil
}

func (curRepo *PostgresCurrencyRepository) Get(ctx context.Context, date time.Time) ([]models.Currency, error) {
	var (
		currency []models.Currency
		query    string
		args     []interface{}
	)

	if !date.IsZero() {
		query = "SELECT name, value, date FROM currencies WHERE date = $1"
		args = append(args, date)
	} else {
		query = "SELECT name, value, date FROM currencies;"
	}

	err := curRepo.conn.SelectContext(ctx, &currency, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from DB: %w", err)
	}
	return currency, nil
}

func (curRepo *PostgresCurrencyRepository) Update(ctx context.Context, rates []models.Currency) error {
	query := `
        INSERT INTO currencies (name, value, date)
        VALUES ($1, $2, $3)
        ON CONFLICT (name, date) DO NOTHING;
    `

	tx, err := curRepo.conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, rate := range rates {
		_, err = stmt.ExecContext(ctx, rate.Name, rate.Value, rate.Date)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute insert statement: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
