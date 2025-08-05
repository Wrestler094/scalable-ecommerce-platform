package dto

// ====== Register ======

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterResponse struct {
	UserID int64 `json:"user_id"`
}

// ====== Login ======

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

// ====== Refresh ======

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

// ====== Logout ======

type LogoutResponse struct {
	Success bool `json:"success"`
}
