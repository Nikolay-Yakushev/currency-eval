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
	logger     *zap.Logger
	repository repository.DatabaseRepository
}

func NewCurrencyUseCase(logger *zap.Logger, repository repository.DatabaseRepository) (*UseCase, error) {
	return &UseCase{
		logger:     logger,
		repository: repository,
	}, nil
}

func (uc *UseCase) GetCurrentExchangeRateByDate(ctx context.Context, d dto.RequestCurrencyByDateDTO) (*dto.ResponseCurrencyByDateDTO, error) {
	models, err := uc.repository.Get(ctx, d.EffectiveDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data in useCase level. Reason: %v", err)
	}

	if len(*models) == 0 {
		return nil, fmt.Errorf("no currency data available for the specified date")
	}

	var (
		usdValue          float64 = 1 // assuming it always equal to 1
		baseCurrencyValue float64
		UpdatedAt         time.Time
	)

	for _, currency := range *models {
		if UpdatedAt.IsZero() {
			UpdatedAt = currency.Date
		}
		if currency.Name == d.BaseCurrency {
			baseCurrencyValue = currency.Value
			break
		}
	}

	if baseCurrencyValue == 0 {
		return nil, fmt.Errorf("base currency %s not found in the data", d.BaseCurrency)
	}

	if d.BaseCurrency != "USD" {
		baseCurrencyValue /= usdValue
	}

	calculatedCurrencies := make(map[string]float64)
	for _, currency := range *models {
		convertedValue := (currency.Value / usdValue) / baseCurrencyValue
		calculatedCurrencies[currency.Name] = convertedValue
	}
	r := dto.ResponseCurrencyByDateDTO{
		UpdatedAt:         UpdatedAt,
		BaseCurrency:      d.BaseCurrency,
		BaseCurrencyValue: float64(1),
		Currencies:        calculatedCurrencies,
	}
	return &r, nil
}

func (uc *UseCase) GetExchangePairRate(ctx context.Context, d dto.RequestCurrencyPairDTO) (*dto.ResponseCurrencyPairDTO, error) {
	models, err := uc.repository.Get(ctx, time.Time{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data in useCase level. Reason: %v", err)
	}
	if len(*models) == 0 {
		return nil, fmt.Errorf("no currency data available")
	}

	var (
		usdValue            float64 = 1 // assuming it always equal to 1
		updatedAt           time.Time
		baseCurrencyValue   float64
		targetCurrencyValue float64
	)

	for _, currency := range *models {
		if updatedAt.IsZero() {
			updatedAt = currency.Date
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
		return nil, fmt.Errorf("no data available for the specified base currency")
	}
	if targetCurrencyValue == 0 {
		return nil, fmt.Errorf("no data available for the specified target currency")
	}

	if d.BaseCurrency != "USD" {
		baseCurrencyValue /= usdValue
	}
	if d.BaseCurrency != "USD" {
		targetCurrencyValue = (targetCurrencyValue / usdValue) / baseCurrencyValue
	}

	response := &dto.ResponseCurrencyPairDTO{
		BaseCurrency:        d.BaseCurrency,
		BaseCurrencyValue:   1,
		TargetCurrency:      d.TargetCurrency,
		TargetCurrencyValue: targetCurrencyValue,
		UpdateAt:            updatedAt,
	}

	return response, nil
}
