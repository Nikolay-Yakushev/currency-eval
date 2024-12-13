package currency

import (
	"context"
	"currency_eval/internal/dto"
	"currency_eval/internal/models"
	"currency_eval/internal/usecase/currency/mock"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"math"
	"testing"
	"time"
)

var apiKey = "test_api_key"

func genTestData(date time.Time) []models.Currency {
	mockCurrencies := map[string]float64{
		"USD": 1.0,    // United States Dollar
		"EUR": 0.85,   // Euro
		"JPY": 110.0,  // Japanese Yen
		"GBP": 0.75,   // British Pound
		"AUD": 1.35,   // Australian Dollar
		"CAD": 1.25,   // Canadian Dollar
		"CHF": 0.92,   // Swiss Franc
		"CNY": 6.45,   // Chinese Yuan
		"INR": 73.0,   // Indian Rupee
		"BRL": 5.25,   // Brazilian Real
		"RUB": 74.0,   // Russian Ruble
		"ZAR": 15.0,   // South African Rand
		"SEK": 8.7,    // Swedish Krona
		"NOK": 8.5,    // Norwegian Krone
		"MXN": 20.0,   // Mexican Peso
		"SGD": 1.35,   // Singapore Dollar
		"HKD": 7.8,    // Hong Kong Dollar
		"KRW": 1150.0, // South Korean Won
		"NZD": 1.4,    // New Zealand Dollar
		"TRY": 9.0,    // Turkish Lira
	}
	var fixtures []models.Currency
	for name, value := range mockCurrencies {
		fixtures = append(fixtures, models.Currency{
			Name:  name,
			Value: value,
			Date:  date,
		})
	}
	return fixtures
}

func floatsAreEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestGetCurrentExchangeRateByDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	mockRepo := mock.NewMockCurrencyRepository(ctrl)

	mockModels := genTestData(mockDate)

	mockRepo.EXPECT().Get(context.Background(), mockDate).Return(mockModels, nil).Times(3)
	mockRepo.EXPECT().Get(context.Background(), mockDate).Return(nil, errors.New("database error")).Times(1)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	uc, err := NewCurrencyUseCase(logger, mockRepo, apiKey)
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		name        string
		request     dto.UseCaseRequestCurrencyByDateDTO
		expected    *dto.UseCaseResponseCurrencyByDateDTO
		expectError bool
	}{
		{
			name: "Base USD",
			request: dto.UseCaseRequestCurrencyByDateDTO{
				BaseCurrency:  "USD",
				EffectiveDate: mockDate,
			},
			expected: &dto.UseCaseResponseCurrencyByDateDTO{
				UpdatedAt:         mockDate,
				BaseCurrency:      "USD",
				BaseCurrencyValue: 1.0,
				Currencies: map[string]float64{
					"EUR": 0.85,
					"JPY": 110.0,
				},
			},
			expectError: false,
		},
		{
			name: "Base EUR",
			request: dto.UseCaseRequestCurrencyByDateDTO{
				BaseCurrency:  "EUR",
				EffectiveDate: mockDate,
			},
			expected: &dto.UseCaseResponseCurrencyByDateDTO{
				UpdatedAt:         mockDate,
				BaseCurrency:      "EUR",
				BaseCurrencyValue: 1.0,
				Currencies: map[string]float64{
					"USD": 1.18,
					"JPY": 129.41,
				},
			},
			expectError: false,
		},
		{
			name: "Invalid Base Currency",
			request: dto.UseCaseRequestCurrencyByDateDTO{
				BaseCurrency:  "INVALID",
				EffectiveDate: mockDate,
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Database Error",
			request: dto.UseCaseRequestCurrencyByDateDTO{
				BaseCurrency:  "USD",
				EffectiveDate: mockDate,
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := uc.GetCurrentExchangeRateByDate(context.Background(), tc.request)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.BaseCurrency, response.BaseCurrency)
				assert.Equal(t, tc.expected.BaseCurrencyValue, response.BaseCurrencyValue)
				for currency, expectedValue := range tc.expected.Currencies {
					assert.Contains(t, response.Currencies, currency)
					actualValue := response.Currencies[currency]
					assert.True(t, floatsAreEqual(expectedValue, actualValue, 1e-2),
						"Expected value %.6f, got %.6f for currency %s", expectedValue, actualValue, currency)
				}
			}
		})
	}
}

func TestGetExchangePairRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	mockRepo := mock.NewMockCurrencyRepository(ctrl)

	mockModels := genTestData(mockDate)
	mockRepo.EXPECT().Get(context.Background(), time.Time{}).Return(mockModels, nil).Times(1)
	mockRepo.EXPECT().Get(context.Background(), time.Time{}).Return(nil, errors.New("database error")).Times(1)
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	uc, err := NewCurrencyUseCase(logger, mockRepo, apiKey)
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		name        string
		request     dto.UseCaseRequestCurrencyPairDTO
		expected    *dto.UseCaseResponseCurrencyPairDTO
		expectError bool
	}{
		{
			name: "Base USD",
			request: dto.UseCaseRequestCurrencyPairDTO{
				BaseCurrency:   "USD",
				TargetCurrency: "EUR",
			},
			expected: &dto.UseCaseResponseCurrencyPairDTO{
				BaseCurrency:        "USD",
				BaseCurrencyValue:   1.0,
				TargetCurrency:      "EUR",
				TargetCurrencyValue: 0.85,
				UpdateAt:            mockDate,
			},
			expectError: false,
		},
		{
			name: "Shit Base Currency",
			request: dto.UseCaseRequestCurrencyPairDTO{
				BaseCurrency:   "Shit",
				TargetCurrency: "EUR",
			},
			expected:    &dto.UseCaseResponseCurrencyPairDTO{},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := uc.GetExchangePairRate(context.Background(), tc.request)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.BaseCurrency, response.BaseCurrency)
				assert.Equal(t, tc.expected.BaseCurrencyValue, response.BaseCurrencyValue)
				assert.Equal(t, tc.expected.TargetCurrencyValue, response.TargetCurrencyValue)
				assert.Equal(t, tc.expected.UpdateAt, response.UpdateAt)
			}
		})
	}
}
