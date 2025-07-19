package user

type UserRole string

const (
	RolePlatformAdmin UserRole = "PLATFORM_ADMIN"
	RoleTenantAdmin   UserRole = "TENANT_ADMIN"
	RoleUser          UserRole = "USER"
)

type CreateUserRequest struct {
	TenantID string
	Email    string
	Phone    string
	Password string
	IsActive bool
	Role     UserRole // PLATFORM_ADMIN, TENANT_ADMIN, USER
}
