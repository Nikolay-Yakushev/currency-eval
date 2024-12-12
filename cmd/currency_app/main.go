package main

import (
	"context"
	"currency_eval/internal/delivery/http"
	_ "currency_eval/internal/docs"
	appLogger "currency_eval/internal/pkg/logger"
	"currency_eval/internal/repository/postgres"
	postgresCurrencyRepo "currency_eval/internal/repository/postgres/currency"
	currencyUsecase "currency_eval/internal/usecase/currency"
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
	log.Println("current working directory", zap.String("path", absPath))

	conf, err := NewConf(absPath)
	if err != nil {
		log.Fatalf("failed to launch currency_app config %v", err)
	}
	logger, err := appLogger.NewLogger(conf.LogLevel)
	if err != nil {
		log.Fatalf("Failed to launch currency_app logger %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("failed to Sync logger %v", err)
		}
	}()

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
		MaxOpenConns:     100,
		MaxIdleConns:     10,
		ConnMaxLifetime:  time.Second * 30,
	}

	currencyRepository, err := postgresCurrencyRepo.NewCurrencyRepository(logger.Named("postgresRepo"), pgConf)
	if err != nil {
		logger.Fatal("failed to launch currency repository", zap.Error(err))
	}
	currencyUseCase, err := currencyUsecase.NewCurrencyUseCase(
		logger.Named("currencyUC"), currencyRepository, conf.CurrencyServiceApiKey,
	)
	if err != nil {
		logger.Fatal("failed to launch currency usecase", zap.Error(err))
	}
	htpConf := http.Config{
		RestApiPort: conf.RestApiPort,
	}

	HTTPController, err := http.NewController(ctx, logger, htpConf, currencyUseCase)
	if err != nil {
		logger.Fatal("failed to initialize currency_app http controller", zap.Error(err))
	}

	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()

	if err := HTTPController.Start(); err != nil {
		logger.Fatal("failed to launch currency_app http_controller", zap.Error(err))
	}

	<-ctx.Done()
	if err := HTTPController.Stop(ctx); err != nil {
		logger.Error("failed to stop currency_app gracefully", zap.Error(err))
	}
	currencyUseCase.Stop()
}
