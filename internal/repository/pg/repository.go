package pg

import (
	"context"
	"currency_eval/internal/models"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"time"
)

func (pgr *PostgresRepository) Get(ctx context.Context, date time.Time) (*[]models.Currency, error) {
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

	err := pgr.conn.SelectContext(ctx, &currency, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from DB: %w", err)
	}
	return &currency, nil
}

func (pgr *PostgresRepository) Update(ctx context.Context, rates *[]models.Currency) error {
	query := `
        INSERT INTO currencies (name, value, date)
        VALUES ($1, $2, $3)
        ON CONFLICT (name, date) DO NOTHING;
    `

	tx, err := pgr.conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, rate := range *rates {
		_, err = stmt.ExecContext(ctx, rate.Name, rate.Value, rate.Date)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute insert statement: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
