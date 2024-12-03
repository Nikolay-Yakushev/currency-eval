package http

import (
	"currency_eval/internal/dto"
	"github.com/gofiber/fiber/v2"
)

// @Tags Currencies
// @Accept json
// @Produce json
// @Param request body dto.RequestCurrencyPairDTO true "Request body"
// @Success 200 {object} dto.ResponseCurrencyPairDTO
// @Router /currencies [post]
func (c *Controller) Pair(ctx *fiber.Ctx) error {
	var request dto.RequestCurrencyPairDTO

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}
	res, err := c.CurrencyUc.GetExchangePairRate(c.ctx, request.ToUpperCase())
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

// @Tags Currencies
// @Accept json
// @Produce json
// @Param request body dto.RequestCurrencyByDateDTO true "Request body"
// @Success 200 {object} dto.ResponseCurrencyByDateDTO
// @Router /currencies_with_date [post]
func (c *Controller) DatePair(ctx *fiber.Ctx) error {
	var request dto.RequestCurrencyByDateDTO

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}
	res, err := c.CurrencyUc.GetCurrentExchangeRateByDate(c.ctx, request.ToUpperCase())
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
