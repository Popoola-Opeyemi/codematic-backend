package handler

import (
	"codematic/internal/domain/auth"
	"codematic/internal/domain/user"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

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
	authGroup.Post("/nonce", a.GetNonce)
	authGroup.Post("/wallet", a.WalletLogin)

	// Protected routes group with JWT middleware
	protected := authGroup.Use(middleware.JWTMiddleware(
		env.JWTManager,
		env.CacheManager,
	))
	protected.Post("/logout", a.Logout)
	protected.Post("/refresh", a.RefreshToken)

	return nil

}

func (a *Auth) Login(c *fiber.Ctx) error {
	var req auth.LoginRequest

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := c.BodyParser(&req); err != nil {
		a.env.Logger.Error("Failed to Parse Body", zap.Error(err))
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	sessionInfo := model.UserSessionInfo{
		UserAgent: c.Get("User-Agent"),
		IPAddress: c.IP(),
		TokenID:   uuid.New().String(),
	}

	auth, err := a.service.Login(ctx, &req, &sessionInfo)
	if err != nil {
		a.env.Logger.Error("Failed to loginn", zap.Error(err))
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

func (a *Auth) RefreshToken(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req := auth.RefreshTokenRequest{}
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	auth, err := a.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccessResponse(c, 200, auth)

}

func (a *Auth) GetNonce(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	walletAddress := c.Query("wallet_address")
	if walletAddress == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "wallet_address is required")
	}

	nonce, err := a.service.GetNonce(ctx, walletAddress)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to Get nonce")
	}

	return utils.SendSuccessResponse(c, 200, nonce)

}

func (a *Auth) WalletLogin(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req := auth.WalletLoginRequest{}

	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if req.WalletAddress == "" || req.Signature == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "wallet_address and signature are required")
	}

	// Prepare session info (user agent, IP, etc)
	sessionInfo := model.UserSessionInfo{
		UserAgent: c.Get("User-Agent"),
		IPAddress: c.IP(),
		TokenID:   uuid.New().String(),
	}

	walletInfo := auth.WalletLoginInfo{
		WalletAddress: req.WalletAddress,
		Nonce:         req.Nonce,
		Signature:     req.Signature,
		Message:       fmt.Sprintf("%s %s", a.env.Config.WALLET_AUTH_MESSAGE, req.Nonce),
	}

	// Call the service layer to handle Wallet Login
	authResult, err := a.service.WalletLogin(ctx, &walletInfo, &sessionInfo)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	// Return success response with auth tokens
	return utils.SendSuccessResponse(c, fiber.StatusOK, authResult)
}
