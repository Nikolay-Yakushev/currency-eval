package dto

import (
	"strings"
	"time"
)

// Controller
type ControllerRequestCurrencyPair struct {
	BaseCurrency   string
	TargetCurrency string
}

type ControllerResponseCurrencyPair struct {
	BaseCurrency        string
	BaseCurrencyValue   float64
	TargetCurrency      string
	TargetCurrencyValue float64
	UpdateAt            time.Time
}

type ControllerRequestCurrencyByDateDTO struct {
	BaseCurrency  string
	EffectiveDate time.Time
}

type ControllerResponseCurrencyByDateDTO struct {
	UpdatedAt         time.Time
	BaseCurrency      string
	BaseCurrencyValue float64
	Currencies        map[string]float64 // {EUR: 1.23} // value relative to BaseCurrency value
}

// Usecase
type UseCaseRequestCurrencyPairDTO struct {
	BaseCurrency   string //relative to which currency rates should be calculated
	TargetCurrency string
}

func (dto *UseCaseRequestCurrencyPairDTO) ToUpperCase() UseCaseRequestCurrencyPairDTO {
	dto.BaseCurrency = strings.ToUpper(dto.BaseCurrency)
	dto.TargetCurrency = strings.ToUpper(dto.TargetCurrency)
	return *dto
}

type UseCaseResponseCurrencyPairDTO struct {
	BaseCurrency        string
	BaseCurrencyValue   float64
	TargetCurrency      string
	TargetCurrencyValue float64
	UpdateAt            time.Time
}

type UseCaseRequestCurrencyByDateDTO struct {
	BaseCurrency  string
	EffectiveDate time.Time
}

func (dto *UseCaseRequestCurrencyByDateDTO) ToUpperCase() UseCaseRequestCurrencyByDateDTO {
	dto.BaseCurrency = strings.ToUpper(dto.BaseCurrency)
	return *dto
}

type UseCaseResponseCurrencyByDateDTO struct {
	UpdatedAt         time.Time
	BaseCurrency      string
	BaseCurrencyValue float64
	Currencies        map[string]float64 // {EUR: 1.23} // value relative to BaseCurrency value
}
