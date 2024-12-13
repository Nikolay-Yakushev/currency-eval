package currency

import (
	"context"
	"currency_eval/internal/dto"
	"currency_eval/internal/models"
	"currency_eval/internal/repository"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

type UseCase struct {
	logger             *zap.Logger
	c                  *cron.Cron
	currencyRepository repository.CurrencyRepository
}

func (uc *UseCase) Stop() {
	uc.c.Stop()
}

func NewCurrencyUseCase(logger *zap.Logger, repository repository.CurrencyRepository, apiKey string) (*UseCase, error) {
	uc := UseCase{
		logger:             logger,
		currencyRepository: repository,
	}

	c := cron.New()
	_, err := c.AddFunc("@every 12h", func() {
		ctx := context.Background()
		err := FetchCurrentCurrencies(ctx, logger, apiKey, uc)
		if err != nil {
			logger.Error("failed to fetch currencies")
		} else {
			logger.Info("fetch currencies completed successfully")
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to schedule cron job: %v", err)
	}
	c.Start()
	logger.Info("cron initialized")
	uc.c = c
	return &uc, nil
}

func (uc *UseCase) GetCurrentExchangeRateByDate(ctx context.Context, d dto.UseCaseRequestCurrencyByDateDTO) (dto.UseCaseResponseCurrencyByDateDTO, error) {
	uc.logger.Info("started fetching")
	var resp = dto.UseCaseResponseCurrencyByDateDTO{}

	currencies, err := uc.currencyRepository.Get(ctx, d.EffectiveDate)
	if err != nil {
		return resp, fmt.Errorf("failed to fetch data: %w", err)
	}

	if len(currencies) == 0 {
		return resp, fmt.Errorf("no currency data available for the specified date")
	}

	var (
		baseCurrencyValue float64
		UpdatedAt         time.Time
	)

	for _, currency := range currencies {
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
	for _, currency := range currencies {
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

	currencies, err := uc.currencyRepository.Get(ctx, time.Time{})
	if err != nil {
		return resp, fmt.Errorf("failed to fetch data in useCase level. Reason: %w", err)
	}
	if len(currencies) == 0 {
		return resp, fmt.Errorf("no currency data available")
	}

	var (
		updatedAt           time.Time
		baseCurrencyValue   float64
		targetCurrencyValue float64
	)

	for _, currency := range currencies {
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

type Result struct {
	Date  string            `json:"date"`
	Base  string            `json:"base"`
	Rates map[string]string `json:"rates"`
}

func FetchCurrentCurrencies(ctx context.Context, logger *zap.Logger, apiKey string, uc UseCase) error {
	logger.Debug("starting update db function")
	requestURL := fmt.Sprintf("https://api.currencyfreaks.com/v2.0/rates/latest?apikey=%s", apiKey)
	resp, err := http.Get(requestURL)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %v", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Error("failed to close resBody")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: HTTP status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	var res Result

	err = json.Unmarshal(body, &res)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	const dateFormat = "2006-01-02 15:04:05-07"
	parsedDate, err := time.Parse(dateFormat, res.Date)
	if err != nil {
		return fmt.Errorf("error parsing date: %v", err)
	}

	var currencies []models.Currency
	for name, value := range res.Rates {
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse input data %v", err)
		}
		c := models.Currency{
			Name:  name,
			Value: value,
			Date:  parsedDate,
		}
		currencies = append(currencies, c)
	}
	baseCurrency := models.Currency{
		Name:  res.Base,
		Value: float64(1),
		Date:  parsedDate,
	}
	currencies = append(currencies, baseCurrency)
	if err := uc.currencyRepository.Update(ctx, currencies); err != nil {
		return fmt.Errorf("failed to insert data. Reason %v", err)
	}
	return nil
}
