package middleware

import (
	"codematic/internal/config"
	"codematic/internal/domain/idempotency"
	"codematic/internal/shared/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IdempotencyMiddleware struct {
	Repo idempotency.Repository
}

func NewIdempotencyMiddleware(repo idempotency.Repository) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{Repo: repo}
}

func (m *IdempotencyMiddleware) Handle(c *fiber.Ctx) error {
	logger := config.GetLogger().Sugar()
	logger.Info("Idempotency Middleware Init")

	key := c.Get("Idempotency-Key")
	if key == "" {
		logger.Error("Missing Idempotency-Key header")
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Missing Idempotency-Key header")
	}

	tenantID := utils.ExtractTenantFromJWT(c)
	if tenantID == "" {
		logger.Error("Missing X-Tenant-ID header")
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Missing X-Tenant-ID header")
	}

	if _, err := uuid.Parse(tenantID); err != nil {
		logger.Errorf("Invalid Tenant ID UUID: %v", err)
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid X-Tenant-ID header format")
	}

	endpoint := c.OriginalURL()
	requestHash := utils.HashRequestBody(c.Body())

	logger.Infof("Idempotency ➤ tenantID=%s | key=%s | endpoint=%s | requestHash=%s", tenantID, key, endpoint, requestHash)

	// Check for cached record
	record, err := m.Repo.Get(c.Context(), tenantID, key, endpoint, requestHash)
	if err == nil && record != nil && record.ID.Valid && record.RequestHash == requestHash {
		logger.Info("Returning cached idempotent response")

		if len(record.ResponseBody) > 0 {
			c.Response().Header.Set("Content-Type", fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(int(record.StatusCode.Int32)).Send(record.ResponseBody)
		}

		// Unexpected: record exists but response body is empty
		logger.Error("Cached response body is empty")
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Cached response is invalid")
	}

	if err != nil {
		logger.Errorf("Error checking idempotency record: %v", err)
	} else {
		logger.Info("No idempotency record found")
	}

	// Proceed with request
	logger.Info("Calling next handler...")
	if err := c.Next(); err != nil {
		logger.Errorf("Handler error: %v", err)
		return err
	}

	status := c.Response().StatusCode()
	responseBody := c.Response().Body()

	// Validate response is JSON
	if !json.Valid(responseBody) {
		logger.Error("Handler returned invalid JSON")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid response format",
		})
	}

	cloned := make([]byte, len(responseBody))
	copy(cloned, responseBody)

	record, err = m.Repo.Get(c.Context(), tenantID, key, endpoint, requestHash)
	if err != nil || record == nil || !record.ID.Valid {
		logger.Info("Saving new idempotency record")

		userID := c.Locals("user_id")
		var userIDStr string
		if userID != nil {
			userIDStr, _ = userID.(string)
		}

		id := uuid.NewString()
		err := m.Repo.Create(c.Context(), idempotency.CreateParams{
			ID:           id,
			TenantID:     tenantID,
			UserID:       userIDStr,
			Key:          key,
			Endpoint:     endpoint,
			RequestHash:  requestHash,
			ResponseBody: cloned,
			StatusCode:   status,
		})
		if err != nil {
			logger.Errorf("Failed to create idempotency record: %v", err)
		}
	} else {
		logger.Info("Updating existing idempotency record")
		_, _ = m.Repo.UpdateResponse(c.Context(), idempotency.UpdateResponseParams{
			TenantID:     tenantID,
			Key:          key,
			Endpoint:     endpoint,
			ResponseBody: cloned,
			StatusCode:   status,
		})
	}

	return nil
}
