package model

type UserRole string

const (
	RolePlatformAdmin UserRole = "PLATFORM_ADMIN"
	RoleTenantAdmin   UserRole = "TENANT_ADMIN"
	RoleUser          UserRole = "USER"
)

// String returns the string representation of the UserRole
func (r UserRole) String() string {
	return string(r)
}
