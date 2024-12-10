package http

import (
	"currency_eval/internal/dto"
	"github.com/gofiber/fiber/v2"
)

// @Tags Currencies
// @Accept json
// @Produce json
// @Param request body dto.ControllerRequestCurrencyPair  true "Request body"
// @Success 200 {object} dto.ControllerResponseCurrencyPair
// @Router /currencies [post]
func (c *Controller) Pair(ctx *fiber.Ctx) error {
	var req dto.ControllerRequestCurrencyPair

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	ucReq := dto.UseCaseRequestCurrencyPairDTO(req)
	res, err := c.CurrencyUc.GetExchangePairRate(c.ctx, ucReq.ToUpperCase())
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}
	controllerResp := dto.ControllerResponseCurrencyPair(res)

	return ctx.Status(fiber.StatusOK).JSON(controllerResp)
}

// @Tags Currencies
// @Accept json
// @Produce json
// @Param request body dto.ControllerRequestCurrencyByDateDTO true "Request body"
// @Success 200 {object} dto.ControllerResponseCurrencyByDateDTO
// @Router /currencies_with_date [post]
func (c *Controller) DatePair(ctx *fiber.Ctx) error {
	var r dto.ControllerRequestCurrencyByDateDTO

	if err := ctx.BodyParser(&r); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}
	ucReq := dto.UseCaseRequestCurrencyByDateDTO(r)

	res, err := c.CurrencyUc.GetCurrentExchangeRateByDate(c.ctx, ucReq.ToUpperCase())
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}
	controllerRes := dto.ControllerResponseCurrencyByDateDTO(res)

	return ctx.Status(fiber.StatusOK).JSON(controllerRes)
}
