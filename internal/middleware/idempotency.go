package middleware

import (
	"bytes"
	db "codematic/internal/infrastructure/db/sqlc"
	"codematic/internal/shared/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type IdempotencyMiddleware struct {
	DB *db.Queries
}

func NewIdempotencyMiddleware(db *db.Queries) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{DB: db}
}

func (m *IdempotencyMiddleware) Handle(c *fiber.Ctx) error {
	key := c.Get("Idempotency-Key")
	if key == "" {

		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "missing Idempotency-Key header")
	}

	tenantID := c.Get("X-Tenant-ID")
	if tenantID == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "missing X-Tenant-ID header")
	}

	endpoint := c.OriginalURL()
	requestHash := utils.HashRequestBody(c.Body())

	tid, err := utils.StringToPgUUID(tenantID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "could not parse tenant ID")
	}

	// Check for existing idempotency record
	record, err := m.DB.GetIdempotencyRecord(c.Context(), db.GetIdempotencyRecordParams{
		TenantID:       tid,
		IdempotencyKey: key,
		Endpoint:       endpoint,
		RequestHash:    requestHash,
	})
	if err == nil && record.ID.Valid {
		// Found: return saved response
		var resp map[string]interface{}
		_ = json.Unmarshal(record.ResponseBody, &resp)
		return c.Status(int(record.StatusCode.Int32)).JSON(resp)
	}

	var buf bytes.Buffer

	c.Response().SetBodyStream(&buf, -1)
	if err := c.Next(); err != nil {
		return err
	}

	status := c.Response().StatusCode()

	body := buf.Bytes()
	if len(body) == 0 {
		body = c.Response().Body()
	}

	id := uuid.New()

	pgID, _ := utils.StringToPgUUID(id.String())
	_, _ = m.DB.SaveIdempotencyRecord(c.Context(), db.SaveIdempotencyRecordParams{
		ID:             pgID,
		TenantID:       tid,
		UserID:         pgtype.UUID{},
		IdempotencyKey: key,
		Endpoint:       endpoint,
		RequestHash:    requestHash,
		ResponseBody:   body,
		StatusCode:     pgtype.Int4{Int32: int32(status), Valid: true},
	})
	return nil
}
