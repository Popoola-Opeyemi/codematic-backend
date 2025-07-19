package middleware

import (
	"codematic/internal/domain/idempotency"
	"codematic/internal/shared/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type IdempotencyMiddleware struct {
	Repo idempotency.Repository
}

func NewIdempotencyMiddleware(repo idempotency.Repository) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{Repo: repo}
}

func (m *IdempotencyMiddleware) Handle(c *fiber.Ctx) error {
	key := c.Get("Idempotency-Key")
	if key == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Missing Idempotency-Key header")
	}

	tenantID := c.Get("X-Tenant-ID")
	if tenantID == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Missing X-Tenant-ID header")
	}

	endpoint := c.OriginalURL()
	requestHash := utils.HashRequestBody(c.Body())

	// Try to find an existing idempotency record
	record, err := m.Repo.Get(c.Context(), tenantID, key, endpoint, requestHash)
	if err == nil && record != nil && record.ID.Valid && record.RequestHash == requestHash {
		var resp map[string]interface{}
		_ = json.Unmarshal(record.ResponseBody, &resp)
		return c.Status(int(record.StatusCode.Int32)).JSON(resp)
	}

	// Proceed to handler
	if err := c.Next(); err != nil {
		return err
	}

	// After handler logic, cache the response
	status := c.Response().StatusCode()

	var responseMap map[string]interface{}
	_ = json.Unmarshal(c.Response().Body(), &responseMap)

	_, _ = m.Repo.UpdateResponse(c.Context(), idempotency.UpdateResponseParams{
		TenantID:     tenantID,
		Key:          key,
		Endpoint:     endpoint,
		ResponseBody: responseMap,
		StatusCode:   status,
	})

	return nil
}
