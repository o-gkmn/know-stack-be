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
}

type LoginResponse struct {
	Token string `json:"token"`
}

type SetClaimsRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	ClaimIDs []uint `json:"claim_ids" binding:"required"`
}

type SetClaimsResponse struct {
	Message string `json:"message"`
}