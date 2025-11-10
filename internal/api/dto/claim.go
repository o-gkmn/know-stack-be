package dto

type CreateClaimRequest struct {
	Name string `json:"name" binding:"required,min=3,max=30"`
}

type CreateClaimResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UpdateClaimRequest struct {
	ID   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required,min=3,max=30"`
}

type UpdateClaimResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type DeleteClaimRequest struct {
	ID uint `json:"id" binding:"required"`
}

type DeleteClaimResponse struct {
	Message string `json:"message"`
}

type GetClaimsResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}