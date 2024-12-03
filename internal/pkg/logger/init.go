package logger

import (
	"currency_eval/internal/config"
	"fmt"
	"go.uber.org/zap"
)

func NewLogger(config config.Config) (*zap.Logger, error) {
	loggerConfig := zap.NewProductionConfig()

	level, err := zap.ParseAtomicLevel(config.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger. Reason %w", err)
	}

	loggerConfig.Level = level

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger. Reason %w", err)
	}
	return logger, nil
}
