package main

import (
	"context"
	"currency_eval/internal/delivery/http"
	_ "currency_eval/internal/docs"
	appLogger "currency_eval/internal/pkg/logger"
	"currency_eval/internal/repository/postgres"
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

	conf, err := NewConf(".")
	if err != nil {
		log.Fatalf("Failed to launch app config %v", err)
	}
	logger, err := appLogger.NewLogger(conf.LogLevel)
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

	pgConf := postgres.Config{
		PostgresHost:     conf.PostgresHost,
		PostgresUser:     conf.PostgresUser,
		PostgresPassword: conf.PostgresPassword,
		PostgresDB:       conf.PostgresDB,
		PostgresPort:     conf.PostgresPort,
	}

	currencyRepository, err := postgres.NewCurrencyRepository(logger.Named("postgresRepo"), pgConf)
	if err != nil {
		logger.Fatal("failed to launch currency repository", zap.Error(err))
	}
	currencyUseCase, err := currency.NewCurrencyUseCase(logger.Named("currencyUC"), currencyRepository)
	if err != nil {
		logger.Fatal("failed to launch currency usecase", zap.Error(err))
	}
	htpConf := http.Config{
		RestApiPort: conf.RestApiPort,
	}

	HTTPController, err := http.NewController(ctx, logger, htpConf, currencyUseCase)
	if err != nil {
		logger.Fatal("failed to initialize app http controller", zap.Error(err))
	}

	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()

	// Fetch data every 12 hours
	go func() {
		for {
			select {
			case <-ticker.C:
				err := scripts.FetchCurrentCurrencies(ctx, conf.CurrencyServiceApiKey, currencyRepository)
				if err != nil {
					logger.Error("failed", zap.Error(err))
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	if err := HTTPController.Start(); err != nil {
		logger.Fatal("failed to launch app http_controller", zap.Error(err))
	}

	<-ctx.Done()
	HTTPController.Stop(ctx) //nolint
}
