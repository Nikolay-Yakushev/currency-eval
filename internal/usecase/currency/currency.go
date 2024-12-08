package currency

import (
	"context"
	"currency_eval/internal/dto"
	"currency_eval/internal/repository"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type UseCase struct {
	logger             *zap.Logger
	currencyRepository repository.CurrencyRepository
}

func NewCurrencyUseCase(logger *zap.Logger, repository repository.CurrencyRepository) (*UseCase, error) {
	return &UseCase{
		logger:             logger,
		currencyRepository: repository,
	}, nil
}

func (uc *UseCase) GetCurrentExchangeRateByDate(ctx context.Context, d dto.UseCaseRequestCurrencyByDateDTO) (dto.UseCaseResponseCurrencyByDateDTO, error) {
	var resp = dto.UseCaseResponseCurrencyByDateDTO{}

	models, err := uc.currencyRepository.Get(ctx, d.EffectiveDate)
	if err != nil {
		return resp, fmt.Errorf("failed to fetch data: %w", err)
	}

	if len(models) == 0 {
		return resp, fmt.Errorf("no currency data available for the specified date")
	}

	var (
		baseCurrencyValue float64
		UpdatedAt         time.Time
	)

	for _, currency := range models {
		if UpdatedAt.IsZero() {
			UpdatedAt = currency.Date
		}
		if currency.Name == d.BaseCurrency {
			baseCurrencyValue = currency.Value
			break
		}
	}

	if baseCurrencyValue == 0 {
		return resp, fmt.Errorf("base currency %s not found in the data", d.BaseCurrency)
	}
	usdValue := 1.0 // assuming it always equal to 1

	if d.BaseCurrency != "USD" {
		baseCurrencyValue /= usdValue
	}

	calculatedCurrencies := make(map[string]float64)
	for _, currency := range models {
		convertedValue := (currency.Value / usdValue) / baseCurrencyValue
		calculatedCurrencies[currency.Name] = convertedValue
	}
	r := dto.UseCaseResponseCurrencyByDateDTO{
		UpdatedAt:         UpdatedAt,
		BaseCurrency:      d.BaseCurrency,
		BaseCurrencyValue: float64(1),
		Currencies:        calculatedCurrencies,
	}
	return r, nil
}

func (uc *UseCase) GetExchangePairRate(ctx context.Context, d dto.UseCaseRequestCurrencyPairDTO) (dto.UseCaseResponseCurrencyPairDTO, error) {
	var resp = dto.UseCaseResponseCurrencyPairDTO{}

	models, err := uc.currencyRepository.Get(ctx, time.Time{})
	if err != nil {
		return resp, fmt.Errorf("failed to fetch data in useCase level. Reason: %w", err)
	}
	if len(models) == 0 {
		return resp, fmt.Errorf("no currency data available")
	}

	var (
		updatedAt           time.Time
		baseCurrencyValue   float64
		targetCurrencyValue float64
	)

	for _, currency := range models {
		if updatedAt.IsZero() {
			updatedAt = currency.Date
		}

		if d.BaseCurrency == d.TargetCurrency {
			baseCurrencyValue = 1.0
			targetCurrencyValue = 1.0
			break
		}

		switch currency.Name {
		case d.BaseCurrency:
			baseCurrencyValue = currency.Value
		case d.TargetCurrency:
			targetCurrencyValue = currency.Value
		}

		if baseCurrencyValue > 0 && targetCurrencyValue > 0 {
			break
		}
	}

	if baseCurrencyValue == 0 {
		return resp, fmt.Errorf("no data available for the specified base currency")
	}
	if targetCurrencyValue == 0 {
		return resp, fmt.Errorf("no data available for the specified target currency")
	}

	usdValue := 1.0 // assuming it always equal to 1

	if d.BaseCurrency != "USD" {
		baseCurrencyValue /= usdValue
	}
	if d.BaseCurrency != "USD" {
		targetCurrencyValue = (targetCurrencyValue / usdValue) / baseCurrencyValue
	}

	response := dto.UseCaseResponseCurrencyPairDTO{
		BaseCurrency:        d.BaseCurrency,
		BaseCurrencyValue:   1,
		TargetCurrency:      d.TargetCurrency,
		TargetCurrencyValue: targetCurrencyValue,
		UpdateAt:            updatedAt,
	}

	return response, nil
}
