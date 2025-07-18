package user

type CreateUserRequest struct {
	TenantID string
	Email    string
	Phone    string
	Password string
	IsActive bool
}
