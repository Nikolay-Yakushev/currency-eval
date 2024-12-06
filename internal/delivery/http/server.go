package http

import (
	"context"
	"currency_eval/internal/usecase"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Controller struct {
	logger     *zap.Logger
	app        *fiber.App
	cfg        Config
	ctx        context.Context
	CurrencyUc usecase.CurrencyUseCase
}

func NewController(ctx context.Context, l *zap.Logger, config Config, uc usecase.CurrencyUseCase) (*Controller, error) {
	app := fiber.New(fiber.Config{
		Prefork: false,
	})
	return &Controller{logger: l, CurrencyUc: uc, app: app, cfg: config, ctx: ctx}, nil
}

func (c *Controller) Start() error {
	var err error

	c.initRoutes(c.app)
	go func() {
		stringifyPort := fmt.Sprintf(":%d", c.cfg.RestApiPort)
		c.logger.Info("starting fiber app")
		if err = c.app.Listen(stringifyPort); err != nil {
			c.logger.Error("failed to start  HTTPController", zap.Error(err))
		}
	}()
	return nil
}

func (c *Controller) Stop(ctx context.Context) error {
	err := c.app.ShutdownWithContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to stop. Reason: %w", err)
	}
	c.logger.Info("finished up")
	return nil
}
