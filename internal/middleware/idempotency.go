package middleware

import (
	"codematic/internal/config"
	"codematic/internal/domain/idempotency"
	"codematic/internal/shared/utils"
	"database/sql"
	"encoding/json"
	"errors"

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

	key := c.Get("Idempotency-Key")
	if key == "" {
		return c.Next()
	}

	tenantID := utils.ExtractTenantFromJWT(c)
	requestHash := utils.HashRequestBody(c.Body())
	endpoint := c.OriginalURL()

	// check for any existing record for this (tenant, key, endpoint)
	rec, err := m.Repo.GetByKeyAndEndpoint(c.Context(), tenantID, key, endpoint)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Errorf("idempotency lookup error: %v", err)
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "idempotency lookup failed")
	}

	if rec != nil {
		switch {
		case rec.RequestHash == requestHash:
			// same key and same payload return cached response
			c.Response().Header.Set("Content-Type", fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(int(rec.StatusCode.Int32)).Send(rec.ResponseBody)

		default:
			// same key and different payload reject
			return utils.SendErrorResponse(c, fiber.StatusConflict,
				"idempotency key reused with different payload",
			)
		}
	}

	if err := c.Next(); err != nil {
		return err
	}

	// save
	status := c.Response().StatusCode()
	body := c.Response().Body()
	if !json.Valid(body) {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError,
			"invalid JSON from handler",
		)
	}

	userID, _ := c.Locals("user_id").(string)
	_ = m.Repo.Create(c.Context(), idempotency.CreateParams{
		ID:           uuid.NewString(),
		TenantID:     tenantID,
		UserID:       userID,
		Key:          key,
		Endpoint:     endpoint,
		RequestHash:  requestHash,
		ResponseBody: body,
		StatusCode:   status,
	})

	return nil
}
