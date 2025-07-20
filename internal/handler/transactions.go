package handler

import (
	"codematic/internal/domain/transactions"
	"codematic/internal/domain/user"
	"codematic/internal/middleware"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Transactions struct {
	service     transactions.Service
	userService user.Service
	env         *Environment
}

func (h *Transactions) Init(basePath string, env *Environment) error {
	h.env = env
	h.service = env.Services.Transactions
	h.userService = env.Services.User

	group := env.Fiber.Group(basePath + "/transactions")
	protected := group.Use(middleware.JWTMiddleware(
		env.JWTManager,
		env.CacheManager,
	))

	protected.Get("/:id", h.GetTransactionByID)
	protected.Get("/", h.ListTransactions)

	return nil
}

func (h *Transactions) validateUserActive(c *fiber.Ctx) error {

	userID := utils.ExtractUserIDFromJWT(c)
	if userID == "" {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "User ID not found in token")
	}

	ctx := context.Background()

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil || !user.IsActive.Bool {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "User not found or suspended")
	}

	return nil
}

// GetTransactionByID godoc
// @Summary      Get a single transaction
// @Description  Retrieves a single transaction by ID with access control
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        id    path   string  true  "Transaction ID"
// @Success      200   {object}  transactions.Transaction
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /transactions/{id} [get]
func (h *Transactions) GetTransactionByID(c *fiber.Ctx) error {
	if err := h.validateUserActive(c); err != nil {
		return err
	}
	id := c.Params("id")
	userID := utils.ExtractUserIDFromJWT(c)
	role := utils.ExtractUserRoleFromJWT(c)
	tenantID := utils.ExtractTenantFromJWT(c)

	tx, err := h.service.GetTransactionByID(c.Context(), id)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "transaction not found")
	}

	// Access control
	switch role {

	case model.RoleUser.String():
		if tx.WalletID != userID {
			return utils.SendErrorResponse(c, fiber.StatusForbidden, "forbidden")
		}

	case model.RoleTenantAdmin.String():
		if tx.TenantID != tenantID {
			return utils.SendErrorResponse(c, fiber.StatusForbidden, "forbidden")
		}

	case model.RolePlatformAdmin.String():
	default:
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "forbidden")
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, tx)
}

// ListTransactions godoc
// @Summary      List transactions
// @Description  Lists transactions with access control and optional filters
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        status   query   string  false  "Transaction status filter"
// @Param        limit    query   int     false  "Limit"
// @Param        offset   query   int     false  "Offset"
// @Success      200   {object}  map[string][]transactions.Transaction
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Router       /transactions [get]
func (h *Transactions) ListTransactions(c *fiber.Ctx) error {

	if err := h.validateUserActive(c); err != nil {
		return err
	}

	userID := utils.ExtractUserIDFromJWT(c)

	role := utils.ExtractUserRoleFromJWT(c)

	tenantID := utils.ExtractTenantFromJWT(c)

	status := c.Query("status")

	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	var txs []*transactions.Transaction

	var err error

	h.env.Logger.Sugar().Debug("ListTransactions role is ", role)

	switch role {

	case model.RoleUser.String():

		txs, err = h.service.ListTransactionsByUserID(c.Context(), userID, limit, offset)

	case model.RoleTenantAdmin.String():
		if status != "" {
			txs, err = h.service.ListTransactionsByStatus(c.Context(), status, limit, offset)
		} else {
			txs, err = h.service.ListTransactionsByTenantID(c.Context(), tenantID, limit, offset)
		}

	case model.RolePlatformAdmin.String():
		if status != "" {
			txs, err = h.service.ListTransactionsByStatus(c.Context(), status, limit, offset)
		}

		if status == "" {
			txs, err = h.service.ListAllTransactions(c.Context(), limit, offset)
		}
	default:
		return utils.SendErrorResponse(c, fiber.StatusForbidden, "forbidden")
	}

	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, fiber.Map{"transactions": txs})
}
