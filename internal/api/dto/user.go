package dto

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanumunicode,min=3,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type CreateUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Remember bool   `json:"remember" binding:"boolean"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RefreshResponse struct {
	AccessToken string `json:"refreshToken"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RequestPasswordResetResponse struct {
	IsSuccess bool `json:"isSuccess"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutResponse struct {
	IsSuccess bool `json:"isSuccess"`
}

type SetClaimsRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	ClaimIDs []uint `json:"claim_ids" binding:"required"`
}

type SetClaimsResponse struct {
	Message string `json:"message"`
}
