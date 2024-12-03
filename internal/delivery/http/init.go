package http

import (
	"context"
	"currency_eval/internal/config"
	"currency_eval/internal/usecase"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"sync"
)

type Controller struct {
	logger     *zap.Logger
	once       sync.Once
	app        *fiber.App
	cfg        config.Config
	ctx        context.Context
	CurrencyUc usecase.CurrencyUseCase
}

func NewController(ctx context.Context, l *zap.Logger, config config.Config, uc usecase.CurrencyUseCase) (*Controller, error) {
	app := fiber.New(fiber.Config{
		Prefork: false,
	})
	return &Controller{logger: l, CurrencyUc: uc, app: app, cfg: config, ctx: ctx}, nil
}

func (c *Controller) Start() error {
	c.initRoutes(c.app)
	go func() {
		stringifyPort := fmt.Sprintf(":%d", c.cfg.RestApiPort)
		c.logger.Info("Starting Fiber app..")
		if err := c.app.Listen(stringifyPort); err != nil {
			c.logger.Error("Failed to start Fiber app", zap.Error(err))
		}
	}()

	return nil
}

func (c *Controller) Stop(ctx context.Context) error {
	var err error

	c.once.Do(func() {
		err = c.app.ShutdownWithContext(ctx)
	})

	if err != nil {
		return fmt.Errorf("failed to stop. Reason: %w", err)
	}
	c.logger.Info("finished up")
	return nil
}
