package usecase

import (
	"context"
	"currency_eval/internal/dto"
)

type CurrencyUseCase interface {
	GetExchangePairRate(ctx context.Context, dto dto.UseCaseRequestCurrencyPairDTO) (dto.UseCaseResponseCurrencyPairDTO, error)
	GetCurrentExchangeRateByDate(ctx context.Context, dto dto.UseCaseRequestCurrencyByDateDTO) (dto.UseCaseResponseCurrencyByDateDTO, error)
}
