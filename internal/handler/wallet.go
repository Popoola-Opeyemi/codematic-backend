package handler

import (
	"codematic/internal/domain/provider"
	"codematic/internal/domain/wallet"
	"codematic/internal/middleware"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	service wallet.Service
	env     *Environment
}

func (h *Wallet) Init(basePath string, env *Environment) error {
	h.env = env

	walletRepo := wallet.NewRepository(env.DB)

	h.service = wallet.NewService(
		walletRepo,
	)

	group := env.Fiber.Group(basePath + "/wallet")

	// Public webhook route (should not require auth)
	group.Post("/webhook/:provider", h.Webhook)

	// Protected routes group with JWT middleware
	protected := group.Use(middleware.JWTMiddleware(
		env.JWTManager,
		env.CacheManager,
	))

	idm := middleware.NewIdempotencyMiddleware(env.DB.Queries)

	// Add idempotency middleware to transaction-creating routes
	protected.Post("/deposit", idm.Handle, h.Deposit)
	protected.Post("/withdraw", idm.Handle, h.Withdraw)
	protected.Post("/transfer", idm.Handle, h.Transfer)
	protected.Post("/get-balance", h.GetBalance)
	protected.Post("/get-transactions", h.GetTransactions)

	return nil
}

// Deposit godoc
// @Summary      Deposit funds into a wallet
// @Description  Deposits a specified amount into the user's wallet
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        depositRequest  body  object  true  "Deposit request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /wallet/deposit [post]
func (h *Wallet) Deposit(c *fiber.Ctx) error {
	var req wallet.DepositRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c,
			fiber.StatusBadRequest,
			model.ErrInvalidInputError.Error(),
		)
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return utils.SendErrorResponse(c,
			fiber.StatusBadRequest,
			"invalid amount",
		)
	}

	form := wallet.DepositForm{
		UserID:   req.UserID,
		TenantID: req.TenantID,
		WalletID: req.WalletID,
		Amount:   amount,
		Provider: req.Provider,
		Metadata: req.Metadata,
	}

	ctx := context.Background()
	if err := h.service.Deposit(ctx, form); err != nil {
		return utils.SendErrorResponse(c,
			fiber.StatusBadRequest,
			err.Error(),
		)
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, fiber.Map{"status": "success"})

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
	var req wallet.WithdrawalRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c,
			fiber.StatusBadRequest,
			"invalid input",
		)
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid amount"})
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
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
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

	var req wallet.TransferRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid amount"})
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
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success"})
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
	walletID := c.Params("wallet_id")

	ctx := context.Background()

	balance, err := h.service.GetBalance(ctx, walletID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"balance": balance.String()})
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

	walletID := c.Params("wallet_id")

	limit := c.QueryInt("limit", 20)

	offset := c.QueryInt("offset", 0)

	ctx := context.Background()

	txs, err := h.service.GetTransactions(ctx, walletID, limit, offset)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"transactions": txs})
}

// Webhook godoc
// @Summary      Handle provider webhook
// @Description  Processes webhook events from any payment provider
// @Tags         wallet
// @Accept       json
// @Produce      json
// @Param        provider  path  string  true  "Provider code"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /wallet/webhook/{provider} [post]
func (h *Wallet) Webhook(c *fiber.Ctx) error {
	providerCode := c.Params("provider")

	prov, ok := provider.GetProvider(providerCode)
	if !ok {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "unsupported provider")
	}

	err := prov.HandleWebhook(c.Context(), c.Body())
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, fiber.Map{"status": "webhook processed"})
}
