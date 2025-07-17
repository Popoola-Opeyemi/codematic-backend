package middleware

import (
	"codematic/internal/infrastructure/cache"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(jwtManager *utils.JWTManager, cacheManager cache.CacheManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		tokenStr := jwtManager.ExtractTokenFromHeader(authHeader)
		if tokenStr == "" {
			return utils.SendErrorResponse(
				c, fiber.StatusUnauthorized,
				model.ErrMissingOrInvalidAuthorizationHeader.Error(),
			)
		}

		claims, err := jwtManager.ParseJWT(tokenStr)
		if err != nil {
			return utils.SendErrorResponse(
				c, fiber.StatusUnauthorized,
				model.ErrInvalidOrExpiredToken.Error(),
			)
		}

		if cacheManager != nil {
			ctx := context.Background()
			session, err := cacheManager.GetSession(ctx, claims.ID)
			if err != nil || session == nil {
				return utils.SendErrorResponse(
					c, fiber.StatusUnauthorized,
					model.ErrTokenRevoked.Error(),
				)
			}
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("token_id", claims.ID)
		c.Locals("claims", claims)
		return c.Next()
	}
}

// TenantMiddleware extracts tenant from X-Tenant-ID header and sets it in context
func TenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Get("X-Tenant-ID")
		if tenantID == "" {
			return utils.SendErrorResponse(
				c, fiber.StatusBadRequest,
				model.ErrMissingXTenantIDHeader.Error(),
			)
		}
		// Optionally: validate tenantID format (e.g., UUID)
		if len(tenantID) != 36 || strings.Count(tenantID, "-") != 4 {
			return utils.SendErrorResponse(
				c, fiber.StatusBadRequest,
				model.ErrInvalidTenantIDFormat.Error(),
			)
		}
		c.Locals("tenant_id", tenantID)
		return c.Next()
	}
}
