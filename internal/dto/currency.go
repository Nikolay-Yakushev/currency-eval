package dto

import (
	"strings"
	"time"
)

type RequestCurrencyPairDTO struct {
	BaseCurrency   string //relative to which currency rates should be calculated
	TargetCurrency string
}

func (dto *RequestCurrencyPairDTO) ToUpperCase() RequestCurrencyPairDTO {
	dto.BaseCurrency = strings.ToUpper(dto.BaseCurrency)
	dto.TargetCurrency = strings.ToUpper(dto.TargetCurrency)
	return *dto
}

type ResponseCurrencyPairDTO struct {
	BaseCurrency        string
	BaseCurrencyValue   float64
	TargetCurrency      string
	TargetCurrencyValue float64
	UpdateAt            time.Time
}

type RequestCurrencyByDateDTO struct {
	BaseCurrency  string
	EffectiveDate time.Time
}

func (dto *RequestCurrencyByDateDTO) ToUpperCase() RequestCurrencyByDateDTO {
	dto.BaseCurrency = strings.ToUpper(dto.BaseCurrency)
	return *dto
}

type ResponseCurrencyByDateDTO struct {
	UpdatedAt         time.Time
	BaseCurrency      string
	BaseCurrencyValue float64
	Currencies        map[string]float64 // {EUR: 1.23} // value relative to BaseCurrency value
}
