package auth

type (
	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		TenantID string `json:"tenant_id" validate:"required"`
	}

	JwtAuthData struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	LoginResponse struct {
		Auth JwtAuthData `json:"auth"`
		User User        `json:"user"`
	}

	User struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		TenantID  string `json:"tenant_id"`
		Role      string `json:"role"`
	}

	SignupRequest struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Phone     string `json:"phone" validate:"required"`
		Password  string `json:"password" validate:"required,min=8"`
		TenantID  string `json:"tenant_id" validate:"required"`
	}

	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}
)

// {
//   "user": {
//     "id": "user-uuid",
//     "tenant_id": "tenant-uuid",
//     "first_name": "Jane",
//     "last_name": "Doe",
//     "email": "jane.doe@example.com",
//     "username": "janedoe",
//     "phone": "+2348012345678",
//     "is_active": true,
//     "created_at": "2025-07-17T12:00:00Z"
//   },
//   "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//   "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//   "expires_in": 900,
//   "token_type": "Bearer",
//   "setup": {
//     "has_wallet": true,
//     "needs_kyc": true,
//     "onboarding_completed": false
//   }
// }
