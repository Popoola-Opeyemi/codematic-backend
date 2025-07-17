package handler

import (
	"codematic/internal/domain/auth"
	"codematic/internal/domain/user"
	"codematic/internal/middleware"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var validate = validator.New()

type Auth struct {
	service auth.Service
	env     *Environment
}

func (a *Auth) Init(basePath string, env *Environment) error {
	a.env = env

	userRepo := user.NewRepository(env.DB.Queries)
	authRepo := auth.NewRepository(env.DB.Queries)

	a.service = auth.NewService(
		userRepo,
		authRepo,
		env.CacheManager,
		env.JWTManager,
		env.Config,
		env.Logger,
	)

	// Public auth routes
	authGroup := env.Fiber.Group(basePath + "/auth")
	authGroup.Post("/login", a.Login)
	authGroup.Post("/signup", a.Signup)

	// Protected routes group with JWT middleware
	protected := authGroup.Use(middleware.JWTMiddleware(
		env.JWTManager,
		env.CacheManager,
	))
	protected.Post("/logout", a.Logout)
	// protected.Post("/refresh", a.RefreshToken)

	return nil

}

func (a *Auth) Signup(c *fiber.Ctx) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Call service (to be implemented in service layer)
	user, err := a.service.Signup(ctx, &req)
	if err != nil {
		a.env.Logger.Error("Failed to signup", zap.Error(err))
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccessResponse(c, fiber.StatusCreated, user)
}

func (a *Auth) Login(c *fiber.Ctx) error {
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

	auth, err := a.service.Login(ctx, &req, &sessionInfo)
	if err != nil {
		a.env.Logger.Error("Failed to login", zap.Error(err))
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccessResponse(c, 200, auth)
}

func (a *Auth) Logout(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	tokenID, ok := c.Locals("token_id").(string)
	if !ok || tokenID == "" {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	err := a.service.Logout(ctx, tokenID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to logout")
	}

	return utils.SendSuccessResponse(c, 200, "Logged out successfully")
}

// func (a *Auth) RefreshToken(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
// 	defer cancel()

// 	req := auth.RefreshTokenRequest{}
// 	if err := c.BodyParser(&req); err != nil {
// 		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
// 	}

// 	auth, err := a.service.RefreshToken(ctx, req.RefreshToken)
// 	if err != nil {
// 		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
// 	}

// 	return utils.SendSuccessResponse(c, 200, auth)

// }
