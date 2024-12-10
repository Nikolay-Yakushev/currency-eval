package usecase

import (
	"context"
	"currency_eval/internal/dto"
)

type CurrencyUseCase interface {
	GetExchangePairRate(ctx context.Context, d dto.UseCaseRequestCurrencyPairDTO) (dto.UseCaseResponseCurrencyPairDTO, error)
	GetCurrentExchangeRateByDate(ctx context.Context, d dto.UseCaseRequestCurrencyByDateDTO) (dto.UseCaseResponseCurrencyByDateDTO, error)
}
