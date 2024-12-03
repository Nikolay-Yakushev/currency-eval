package scripts

import (
	"context"
	"currency_eval/internal/models"
	"currency_eval/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Result struct {
	Date  string            `json:"date"`
	Base  string            `json:"base"`
	Rates map[string]string `json:"rates"`
}

func FetchCurrentCurrencies(ctx context.Context, apiKey string, repo repository.DatabaseRepository) error {
	requestURL := fmt.Sprintf("https://api.currencyfreaks.com/v2.0/rates/latest?apikey=%s", apiKey)
	resp, err := http.Get(requestURL)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

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
	if err := repo.Update(ctx, &currencies); err != nil {
		return fmt.Errorf("failed to insert data. Reason %v", err)
	}
	return nil
}
