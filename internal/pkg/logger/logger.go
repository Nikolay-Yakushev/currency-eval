package logger

import (
	"fmt"
	"go.uber.org/zap"
)

func NewLogger(logLevel string) (*zap.Logger, error) {
	loggerConfig := zap.NewProductionConfig()

	level, err := zap.ParseAtomicLevel(logLevel)
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
