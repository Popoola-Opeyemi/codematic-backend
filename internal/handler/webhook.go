package handler

import (
	"codematic/internal/domain/webhook"
	"codematic/internal/shared/model"

	"github.com/gofiber/fiber/v2"
)

type Webhook struct {
	service webhook.Service
	env     *Environment
}

func (h *Webhook) Init(basePath string, env *Environment) error {

	h.env = env
	h.service = env.Services.Webhook

	group := env.Fiber.Group(basePath + "/webhook")
	group.Post("/:provider", h.Receive)

	return nil
}

// Receive godoc
// @Summary      Handle provider webhook
// @Description  Processes webhook events from any payment provider
// @Tags         webhook
// @Accept       json
// @Produce      json
// @Param        provider  path  string  true  "Provider code"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /webhook/{provider} [post]
func (h *Webhook) Receive(c *fiber.Ctx) error {
	provider := c.Params("provider")
	payload := c.Body()

	headers := map[string]string{}
	c.Request().Header.VisitAll(func(key, val []byte) {
		headers[string(key)] = string(val)
	})

	h.env.Logger.Sugar().Info("payload", string(payload))
	h.env.Logger.Sugar().Info("headers", headers)

	err := h.service.HandleWebhook(c.Context(), provider, headers, payload)
	if err != nil {
		if err == model.ErrInvalidSignature {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "webhook processed"})
}

// Replay godoc
// @Summary      Replay a webhook event
// @Description  Replays a previously failed webhook event by its ID
// @Tags         webhook
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Webhook Event ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /webhook/replay/{id} [post]
// func (h *Webhook) Replay(c *fiber.Ctx) error {
// 	provider := c.Params("provider")
// 	payload := c.Body()

// 	headers := map[string]string{}
// 	c.Request().Header.VisitAll(func(key, val []byte) {
// 		headers[string(key)] = string(val)
// 	})

// 	err := h.service.ReplayWebhook(c.Context(), provider, headers, payload)
// 	if err != nil {
// 		if err == model.ErrInvalidSignature {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "webhook replayed"})
// }
