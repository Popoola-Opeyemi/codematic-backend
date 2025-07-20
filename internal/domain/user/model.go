package user

import "codematic/internal/shared/model"

type CreateUserRequest struct {
	TenantID string
	Email    string
	Phone    string
	Password string
	IsActive bool
	Role     model.UserRole // PLATFORM_ADMIN, TENANT_ADMIN, USER
}
