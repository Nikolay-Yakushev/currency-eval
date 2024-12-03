package http

import (
	_ "currency_eval/internal/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// Ping endpoint
// @Summary Ping the server
// @Description Responds with a "pong" message
// @Tags Ping
// @Success 200 {object} map[string]string
// @Router /ping [get]
func Ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "pong"})
}

func (c *Controller) initRoutes(a *fiber.App) {
	a.Get("/swagger/*", swagger.HandlerDefault) // Serve Swagger U
	a.Get("/ping", Ping)
	a.Post("/currencies", c.Pair)
	a.Post("/currencies_with_date", c.DatePair)

}
