package middleware

import (
	"codematic/internal/infrastructure/cache"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware enforces authentication for protected routes only (not for login endpoints).
// Use this middleware on routes that require a valid JWT (i.e., after login).
func JWTMiddleware(jwtManager *utils.JWTManager,
	cacheManager cache.CacheManager) fiber.Handler {
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
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

// TenantMiddleware extracts tenant from X-Tenant-ID header and sets it in context.
// Use this only for APIs that expect tenant ID in headers (not used for login, where tenant ID is in the body).
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

// RoleMiddleware enforces that the user has one of the required roles.
// Use this middleware after JWTMiddleware to restrict access to users with specific roles (e.g., TENANT_ADMIN, PLATFORM_ADMIN).
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*model.Claims)
		if !ok || claims == nil {
			return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized: missing claims")
		}
		userRole := claims.Role
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				return c.Next()
			}
		}
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "Forbidden: insufficient role")
	}
}
