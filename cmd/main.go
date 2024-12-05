package main

import (
	"context"
	"currency_eval/internal/config"
	"currency_eval/internal/delivery/http"
	_ "currency_eval/internal/docs"
	"currency_eval/internal/pkg/logger"
	"currency_eval/internal/repository/pg"
	"currency_eval/internal/scripts"
	"currency_eval/internal/usecase/currency"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// @title Currency API
// @version 1.0

func main() {
	absPath, _ := filepath.Abs(".")

	log.Println("Current working directory", zap.String("path", absPath))
	conf, err := config.NewConfig(".")
	if err != nil {
		log.Fatalf("Failed to launch app config %v", err)
	}
	logger, err := logger.NewLogger(conf)
	if err != nil {
		log.Fatalf("Failed to launch app logger %v", err)
	}
	defer logger.Sync() //nolint

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGKILL,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
	)
	defer cancel()

	pgRepository, err := pg.NewPostgresRepository(logger.Named("postgresRepo"), conf)
	if err != nil {
		log.Fatalf("Failed to launch pg repository %v", err)
	}
	currencyUc, err := currency.NewCurrencyUseCase(logger.Named("currencyUC"), pgRepository)
	if err != nil {
		log.Fatalf("Failed to launch uc %v", err)
	}
	httpController, err := http.NewController(ctx, logger, conf, currencyUc)
	if err != nil {
		log.Fatalf("Failed to launch app logger %v", err)
	}

	ticker := time.NewTicker(12 * time.Hour)
	done := make(chan bool, 1)
	defer ticker.Stop()
	// Fetch data every 12 hours
	go func() {
		for {
			select {
			case <-ticker.C:
				err := scripts.FetchCurrentCurrencies(ctx, conf.CurrencyServiceApiKey, pgRepository)
				if err != nil {
					panic(err)
				}
			case <-done:
				return
			}
		}
	}()

	if err := httpController.Start(); err != nil {
		log.Fatalf("Failed to launch app logger %v", err)
	}

	<-ctx.Done()
	done <- true
	httpController.Stop(ctx) //nolint
}
