package usecase

import (
	"context"
	"currency_eval/internal/dto"
)

type CurrencyUseCase interface {
	GetExchangePairRate(ctx context.Context, dto dto.RequestCurrencyPairDTO) (*dto.ResponseCurrencyPairDTO, error)
	GetCurrentExchangeRateByDate(ctx context.Context, dto dto.RequestCurrencyByDateDTO) (*dto.ResponseCurrencyByDateDTO, error)
}
