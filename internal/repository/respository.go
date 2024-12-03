package repository

import (
	"context"
	"currency_eval/internal/models"
	"time"
)

type DatabaseRepository interface {
	Get(ctx context.Context, date time.Time) (*[]models.Currency, error)
	Update(ctx context.Context, rates *[]models.Currency) error
}
