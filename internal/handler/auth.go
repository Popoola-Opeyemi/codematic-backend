package handler

import (
	"codematic/internal/domain/auth"
	"codematic/internal/domain/provider"
	"codematic/internal/domain/tenants"
	"codematic/internal/domain/user"
	"codematic/internal/domain/wallet"
	"codematic/internal/middleware"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Auth struct {
	service auth.Service
	env     *Environment
}

func (h *Auth) Init(basePath string, env *Environment) error {
	h.env = env

	userService := user.NewService(env.DB, env.JWTManager, env.Logger)
	tenantService := tenants.NewService(env.DB, env.JWTManager, env.Logger)
	providerService := provider.NewService(env.DB, env.CacheManager, env.Logger, env.KafkaProducer)

	walletService := wallet.NewService(providerService, env.DB, env.Logger, env.KafkaProducer)

	h.service = auth.NewService(
		env.DB,
		userService,
		walletService,
		tenantService,
		env.CacheManager,
		env.JWTManager,
		env.Config,
		env.Logger,
	)

	// Public auth routes
	authGroup := env.Fiber.Group(basePath + "/auth")
	authGroup.Post("/login", h.Login)
	authGroup.Post("/admin", h.AdminLogin)

	authGroup.Post("/signup", middleware.JWTMiddleware(env.JWTManager, env.CacheManager),
		middleware.RoleMiddleware("TENANT_ADMIN"), h.Signup)

	// Protected routes group with JWT middleware
	protected := authGroup.Use(middleware.JWTMiddleware(
		env.JWTManager,
		env.CacheManager,
	))
	// protected.Get("/users", middleware.RoleMiddleware("TENANT_ADMIN"), h.ListUsers)
	protected.Post("/logout", h.Logout)
	protected.Post("/refresh", h.RefreshToken)

	return nil

}

// Signup godoc
// @Summary      Register a new user
// @Description  Creates a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        signupRequest  body  auth.SignupRequest  true  "Signup request"
// @Success      201  {object}  interface{}
// @Failure      400  {object}  model.ErrorResponse
// @Router       /auth/signup [post]
func (h *Auth) Signup(c *fiber.Ctx) error {
	var req auth.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest,
			model.ErrInvalidInputError.Error())
	}

	if err := validate.Struct(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if len(req.Password) < 8 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, model.ErrPasswordTooShort.Error())
	}

	req.TenantID = utils.ExtractTenantFromJWT(c)
	if req.TenantID == "" {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Tenant ID not found in token")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	user, err := h.service.Signup(ctx, &req)
	if err != nil {
		h.env.Logger.Error("Failed to signup", zap.Error(err))
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccessResponse(c, fiber.StatusCreated, user)
}

// Login godoc
// @Summary      Login a tenant user (regular or tenant admin)
// @Description  Authenticates a tenant user and returns tokens. TenantID must be provided in the request body.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginRequest  body  auth.LoginRequest  true  "Login request"
// @Success      200  {object}  interface{}
// @Failure      400  {object}  model.ErrorResponse
// @Router       /auth/login [post]
func (h *Auth) Login(c *fiber.Ctx) error {
	var req auth.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest,
			model.ErrInvalidInputError.Error())
	}

	if err := validate.Struct(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if req.TenantID == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "tenant_id is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sessionInfo := model.UserSessionInfo{
		UserAgent: c.Get("User-Agent"),
		IPAddress: c.IP(),
		TokenID:   uuid.New().String(),
	}

	authResp, err := h.service.Login(ctx, &req, sessionInfo)
	if err != nil {
		h.env.Logger.Error("Failed to login", zap.Error(err))
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Only allow login for tenant admin or user roles
	user, ok := authResp.(auth.LoginResponse)
	if !ok || (user.User.Role != string(model.RoleTenantAdmin) && user.User.Role != string(model.RoleUser)) {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Only tenant admin or user can login here")
	}

	return utils.SendSuccessResponse(c, 200, authResp)
}

func (h *Auth) Logout(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	tokenID, ok := c.Locals("token_id").(string)
	if !ok || tokenID == "" {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	err := h.service.Logout(ctx, tokenID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to logout")
	}

	return utils.SendSuccessResponse(c, 200, "Logged out successfully")
}

func (h *Auth) RefreshToken(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req := auth.RefreshTokenRequest{}
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	auth, err := h.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccessResponse(c, 200, auth)
}

// AdminLogin godoc
// @Summary      Login a platform admin
// @Description  Authenticates a platform admin and returns tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginRequest  body  auth.LoginRequest  true  "Login request"
// @Success      200  {object}  interface{}
// @Failure      400  {object}  model.ErrorResponse
// @Router       /auth/admin [post]
func (h *Auth) AdminLogin(c *fiber.Ctx) error {
	var req auth.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest,
			model.ErrInvalidInputError.Error())
	}

	if err := validate.Struct(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sessionInfo := model.UserSessionInfo{
		UserAgent: c.Get("User-Agent"),
		IPAddress: c.IP(),
		TokenID:   uuid.New().String(),
	}

	authResp, err := h.service.AdminLogin(ctx, &req, sessionInfo)
	if err != nil {
		h.env.Logger.Error("Failed to login", zap.Error(err))
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Check if user is PLATFORM_ADMIN
	user, ok := authResp.(auth.LoginResponse)
	if !ok || user.User.Role != string(model.RolePlatformAdmin) {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Not an admin user")
	}

	return utils.SendSuccessResponse(c, 200, authResp)
}
