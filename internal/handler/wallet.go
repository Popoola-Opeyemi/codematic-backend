package handler

import (
	"codematic/internal/domain/idempotency"
	"codematic/internal/domain/provider"
	"codematic/internal/domain/user"
	"codematic/internal/domain/wallet"
	"codematic/internal/middleware"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	service     wallet.Service
	userService user.Service
	env         *Environment
}

func (h *Wallet) Init(basePath string, env *Environment) error {
	h.env = env

	idempotencyRepo := idempotency.NewRepository(env.DB.Queries, env.DB.Pool)
	providerService := provider.NewService(env.DB, env.CacheManager, env.Logger, env.KafkaProducer)
	userService := user.NewService(env.DB, env.JWTManager, env.Logger)

	h.service = wallet.NewService(
		env.Logger,
		providerService,
		userService,
		env.DB,
		env.KafkaProducer,
	)
	h.userService = userService

	group := env.Fiber.Group(basePath + "/wallet")

	protected := group.Use(middleware.JWTMiddleware(
		env.JWTManager,
		env.CacheManager,
	))

	idm := middleware.NewIdempotencyMiddleware(idempotencyRepo)

	userOnly := protected.Use(utils.RequireRole(model.RoleUser))

	// Add idempotency middleware to transaction-creating routes
	userOnly.Post("/initiate_deposit", idm.Handle, h.InitiateDeposit)
	userOnly.Post("/withdraw", idm.Handle, h.Withdraw)
	userOnly.Post("/transfer", idm.Handle, h.Transfer)
	userOnly.Post("/get-balance", h.GetBalance)
	userOnly.Post("/get-transactions", h.GetTransactions)

	return nil
}

// validateUserActive checks if the user is active and valid
func (h *Wallet) validateUserActive(c *fiber.Ctx) error {
	userID := utils.ExtractUserIDFromJWT(c)
	if userID == "" {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "User ID not found in token")
	}

	// Check if user exists and is active
	ctx := context.Background()
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	if !user.IsActive.Bool {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "User is suspended")
	}

	return nil
}

func (h *Wallet) InitiateDeposit(c *fiber.Ctx) error {
	if err := h.validateUserActive(c); err != nil {
		return err
	}

	var req wallet.DepositRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c,
			fiber.StatusBadRequest,
			model.ErrInvalidInputError.Error(),
		)
	}

	if !req.Channel.IsValid() {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest,
			"invalid channel")
	}

	if err := validate.Struct(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	userID := utils.ExtractUserIDFromJWT(c)
	if userID == "" {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized,
			"User ID not found in token")
	}

	userEmail := utils.ExtractUserEmailFromJWT(c)
	if userEmail == "" {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized,
			"User Email not found in token")
	}

	tenantID := utils.ExtractTenantFromJWT(c)

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return utils.SendErrorResponse(c,
			fiber.StatusBadRequest,
			"invalid amount",
		)
	}
	metadata := map[string]interface{}{
		"email": userEmail,
	}

	currency := strings.ToUpper(req.Currency)

	form := wallet.DepositForm{
		UserID:   userID,
		TenantID: tenantID,
		Amount:   amount,
		Currency: currency,
		Channel:  string(req.Channel),
		Metadata: metadata,
	}

	ctx := context.Background()

	response, err := h.service.InitiateDeposit(ctx, form)
	if err != nil {
		return utils.SendErrorResponse(c,
			fiber.StatusBadRequest,
			err.Error(),
		)
	}
	h.env.Logger.Sugar().Infof("Deposit response: %+v", response)

	return utils.SendSuccessResponse(c, fiber.StatusOK, fiber.Map{
		"authorization_url": response.AuthorizationURL,
		"reference":         response.Reference,
		"provider":          response.Provider,
		"provider_id":       response.ProviderID,
	})

	// return utils.SendSuccessResponse(c, fiber.StatusOK, response)

}

// Withdraw godoc
// @Summary      Withdraw funds from a wallet
// @Description  Withdraws a specified amount from the user's wallet
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        withdrawRequest  body  object  true  "Withdraw request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /wallet/withdraw [post]
func (h *Wallet) Withdraw(c *fiber.Ctx) error {
	// Validate user is active before proceeding
	if err := h.validateUserActive(c); err != nil {
		return err
	}

	var req wallet.WithdrawalRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c,
			fiber.StatusBadRequest,
			"invalid input",
		)
	}

	// Ensure the user can only withdraw from their own wallet
	userID := utils.ExtractUserIDFromJWT(c)
	if req.UserID != userID {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "You can only withdraw from your own wallet")
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "invalid amount")
	}

	form := wallet.WithdrawalForm{
		UserID:   req.UserID,
		TenantID: req.TenantID,
		WalletID: req.WalletID,
		Amount:   amount,
		Provider: req.Provider,
		Metadata: req.Metadata,
	}

	ctx := context.Background()
	if err := h.service.Withdraw(ctx, form); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, fiber.Map{"status": "success"})

}

// Transfer godoc
// @Summary      Transfer funds between wallets
// @Description  Transfers a specified amount from one wallet to another
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        transferRequest  body  object  true  "Transfer request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /wallet/transfer [post]
func (h *Wallet) Transfer(c *fiber.Ctx) error {
	// Validate user is active before proceeding
	if err := h.validateUserActive(c); err != nil {
		return err
	}

	var req wallet.TransferRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "invalid input")
	}

	// Ensure the user can only transfer from their own wallet
	userID := utils.ExtractUserIDFromJWT(c)
	if req.UserID != userID {
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "You can only transfer from your own wallet")
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "invalid amount")
	}

	form := wallet.TransferForm{
		UserID:       req.UserID,
		TenantID:     req.TenantID,
		FromWalletID: req.FromWalletID,
		ToWalletID:   req.ToWalletID,
		Amount:       amount,
		Metadata:     req.Metadata,
	}

	ctx := context.Background()
	if err := h.service.Transfer(ctx, form); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, fiber.Map{"status": "success"})
}

// GetBalance godoc
// @Summary      Get wallet balance
// @Description  Retrieves the balance of a wallet
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        wallet_id  path  string  true  "Wallet ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /wallet/{wallet_id}/balance [get]
func (h *Wallet) GetBalance(c *fiber.Ctx) error {
	// Validate user is active before proceeding
	if err := h.validateUserActive(c); err != nil {
		return err
	}

	walletID := c.Params("wallet_id")

	// TODO: Add validation to ensure user can only access their own wallet balance
	// This would require checking if the wallet belongs to the authenticated user

	ctx := context.Background()

	balance, err := h.service.GetBalance(ctx, walletID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, fiber.Map{"balance": balance.String()})
}

// GetTransactions godoc
// @Summary      Get wallet transactions
// @Description  Retrieves the transaction history for a wallet
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        wallet_id  path  string  true  "Wallet ID"
// @Param        limit      query int     false "Limit"
// @Param        offset     query int     false "Offset"
// @Success      200  {object}  map[string][]wallet.Transaction
// @Failure      400  {object}  map[string]string
// @Router       /wallet/{wallet_id}/transactions [get]
func (h *Wallet) GetTransactions(c *fiber.Ctx) error {
	// Validate user is active before proceeding
	if err := h.validateUserActive(c); err != nil {
		return err
	}

	walletID := c.Params("wallet_id")

	// TODO: Add validation to ensure user can only access their own wallet transactions
	// This would require checking if the wallet belongs to the authenticated user

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	ctx := context.Background()

	txs, err := h.service.GetTransactions(ctx, walletID, limit, offset)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccessResponse(c, fiber.StatusOK, fiber.Map{"transactions": txs})
}
