package service

import (
	"errors"
	"knowstack/internal/api/dto"
	"knowstack/internal/data/models"
	"knowstack/internal/utils"

	"gorm.io/gorm"
)

var (
	ErrClaimAlreadyExists = errors.New("claim already exists")
	ErrClaimNotFound      = errors.New("claim not found")
)

type ClaimService struct {
	DB *gorm.DB
}

func NewClaimService(db *gorm.DB) *ClaimService {
	return &ClaimService{DB: db}
}

func (s *ClaimService) CreateClaim(req dto.CreateClaimRequest) (*dto.CreateClaimResponse, error) {
	utils.LogInfo("Creating claim: %+v", req)

	if err := s.DB.Where("name = ?", req.Name).First(&models.Claim{}).Error; err == nil {
		utils.LogInfo("Claim already exists: %+v", req.Name)
		return nil, ErrClaimAlreadyExists
	}

	claim := &models.Claim{Name: req.Name}
	if err := s.DB.Create(claim).Error; err != nil {
		utils.LogErrorWithErr("Failed to create claim", err)
		return nil, err
	}

	return &dto.CreateClaimResponse{
		ID:   claim.ID,
		Name: claim.Name,
	}, nil
}

func (s *ClaimService) DeleteClaim(req dto.DeleteClaimRequest) (*dto.DeleteClaimResponse, error) {
	utils.LogInfo("Deleting claim: %+v", req)

	var claim models.Claim
	if err := s.DB.Where("id = ?", req.ID).First(&claim).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogInfo("Claim not found: %+v", req.ID)
			return nil, ErrClaimNotFound
		}
		utils.LogErrorWithErr("Failed to find claim", err)
		return nil, err
	}

	if err := s.DB.Delete(&claim).Error; err != nil {
		utils.LogErrorWithErr("Failed to delete claim", err)
		return nil, err
	}

	return &dto.DeleteClaimResponse{
		Message: "Claim deleted successfully",
	}, nil
}

func (s *ClaimService) GetClaims() ([]dto.GetClaimsResponse, error) {
	utils.LogInfo("Getting claims")

	var claims []models.Claim
	if err := s.DB.Where("deleted_at IS NULL").Find(&claims).Error; err != nil {
		utils.LogErrorWithErr("Failed to get claims", err)
		return nil, err
	}

	response := make([]dto.GetClaimsResponse, len(claims))
	for i, claim := range claims {
		response[i] = dto.GetClaimsResponse{
			ID:   claim.ID,
			Name: claim.Name,
		}
	}

	return response, nil
}

func (s *ClaimService) UpdateClaim(req dto.UpdateClaimRequest) (*dto.UpdateClaimResponse, error) {
	utils.LogInfo("Updating claim: %+v", req)

	var claim models.Claim
	if err := s.DB.Where("id = ?", req.ID).First(&claim).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogInfo("Claim not found: %+v", req.ID)
			return nil, ErrClaimNotFound
		}
		utils.LogErrorWithErr("Failed to find claim", err)
		return nil, err
	}

	if err := s.DB.Where("name = ? AND id != ?", req.Name, req.ID).First(&models.Claim{}).Error; err == nil {
		utils.LogInfo("Claim name already exists: %+v", req.Name)
		return nil, ErrClaimAlreadyExists
	}

	claim.Name = req.Name
	if err := s.DB.Save(&claim).Error; err != nil {
		utils.LogErrorWithErr("Failed to update claim", err)
		return nil, err
	}

	return &dto.UpdateClaimResponse{
		ID:   claim.ID,
		Name: claim.Name,
	}, nil
}
