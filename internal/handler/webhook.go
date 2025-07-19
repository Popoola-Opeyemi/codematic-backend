package handler

import (
	"codematic/internal/domain/webhook"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Webhook handles incoming and replayed webhook events from payment providers.
type Webhook struct {
	service webhook.Service
	env     *Environment
}

// Init sets up the webhook handler routes and injects dependencies.
// @param basePath The base API path (e.g. /api)
// @param env The application environment (dependencies)
// @param service The webhook service implementation
func (h *Webhook) Init(basePath string, env *Environment, service webhook.Service) error {
	h.env = env
	h.service = service

	group := env.Fiber.Group(basePath + "/webhook")
	group.Post("/:provider", h.Receive)
	group.Post("/replay/:id", h.Replay)

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
// @Failure      500  {object}  map[string]string
// @Router       /webhook/{provider} [post]
func (h *Webhook) Receive(c *fiber.Ctx) error {
	provider := c.Params("provider")
	payload := c.Body()
	if err := h.service.ProcessWebhook(c.Context(), provider, payload); err != nil {
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
func (h *Webhook) Replay(c *fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.service.ReplayWebhook(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "webhook replayed"})
}
