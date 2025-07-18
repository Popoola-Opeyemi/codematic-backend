package handler

import (
	"codematic/internal/domain/tenants"
	"codematic/internal/middleware"
	"codematic/internal/shared/model"
	"codematic/internal/shared/utils"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Tenants struct {
	service tenants.Service
	env     *Environment
}

func (h *Tenants) Init(basePath string, env *Environment) error {
	h.env = env

	tenantsRepo := tenants.NewRepository(env.DB.Queries)

	h.service = tenants.NewService(
		tenantsRepo,
		env.JWTManager,
		env.Logger,
	)

	// Public auth routes
	group := env.Fiber.Group(basePath + "/tenant")

	// Protected routes group with JWT middleware
	protected := group.Use(middleware.JWTMiddleware(
		env.JWTManager,
		env.CacheManager,
	))
	protected.Post("/create", h.Create)
	protected.Get("/:id", h.GetByID)
	protected.Get("/slug/:slug", h.GetBySlug)
	protected.Get("/", h.List)
	protected.Put("/:id", h.Update)
	protected.Delete("/:id", h.Delete)

	return nil

}

// Create godoc
// @Summary      Creates a new Tenant
// @Description  Creates a new Tenant
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param tenantsCreateRequest body tenants.CreateTenantRequest true "Tenant creation payload"
// @Success      201  {object}  tenants.Tenant
// @Failure      400  {object}  model.ErrorResponse
// @Router       /tenant/create [post]
func (h *Tenants) Create(c *fiber.Ctx) error {
	var req tenants.CreateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest,
			model.ErrInvalidInputError.Error())
	}

	if err := validate.Struct(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	tenant, err := h.service.CreateTenant(ctx, req)
	if err != nil {
		h.env.Logger.Error("Failed to signup", zap.Error(err))
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccessResponse(c, fiber.StatusCreated, tenant)
}

// GetByID godoc
// @Summary      Get tenant by ID
// @Description  Get tenant by ID
// @Tags         tenants
// @Produce      json
// @Param        id   path      string  true  "Tenant ID"
// @Success      200  {object}  tenants.Tenant
// @Failure      404  {object}  model.ErrorResponse
// @Router       /tenant/{id} [get]
func (h *Tenants) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	tenant, err := h.service.GetTenantByID(ctx, id)
	if err != nil {
		h.env.Logger.Error("Failed to get tenant by ID", zap.Error(err))
		return utils.SendErrorResponse(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccessResponse(c, fiber.StatusOK, tenant)
}

// GetBySlug godoc
// @Summary      Get tenant by slug
// @Description  Get tenant by slug
// @Tags         tenants
// @Produce      json
// @Param        slug   path      string  true  "Tenant Slug"
// @Success      200  {object}  tenants.Tenant
// @Failure      404  {object}  model.ErrorResponse
// @Router       /tenant/slug/{slug} [get]
func (h *Tenants) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	tenant, err := h.service.GetTenantBySlug(ctx, slug)
	if err != nil {
		h.env.Logger.Error("Failed to get tenant by slug", zap.Error(err))
		return utils.SendErrorResponse(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccessResponse(c, fiber.StatusOK, tenant)
}

// List godoc
// @Summary      List tenants
// @Description  List all tenants
// @Tags         tenants
// @Produce      json
// @Success      200  {array}  tenants.Tenant
// @Router       /tenant [get]
func (h *Tenants) List(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	tenants, err := h.service.ListTenants(ctx)
	if err != nil {
		h.env.Logger.Error("Failed to list tenants", zap.Error(err))
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccessResponse(c, fiber.StatusOK, tenants)
}

// Update godoc
// @Summary      Update tenant
// @Description  Update tenant by ID
// @Tags         tenants
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tenant ID"
// @Param        body body      tenants.CreateTenantRequest true "Tenant update payload"
// @Success      200  {object}  tenants.Tenant
// @Failure      400  {object}  model.ErrorResponse
// @Router       /tenant/{id} [put]
func (h *Tenants) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req tenants.CreateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, model.ErrInvalidInputError.Error())
	}
	if err := validate.Struct(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	tenant, err := h.service.UpdateTenant(ctx, id, req.Name, req.Slug)
	if err != nil {
		h.env.Logger.Error("Failed to update tenant", zap.Error(err))
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccessResponse(c, fiber.StatusOK, tenant)
}

// Delete godoc
// @Summary      Delete tenant
// @Description  Delete tenant by ID
// @Tags         tenants
// @Produce      json
// @Param        id   path      string  true  "Tenant ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  model.ErrorResponse
// @Router       /tenant/{id} [delete]
func (h *Tenants) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := h.service.DeleteTenant(ctx, id); err != nil {
		h.env.Logger.Error("Failed to delete tenant", zap.Error(err))
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
