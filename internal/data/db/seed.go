package db

import (
	"knowstack/internal/data/models"
	"knowstack/internal/utils"
)

func SeedData() error {
	// Seed Roles
	roles := []models.Role{
		{Name: "user", IsDefault: true},
		{Name: "admin", IsDefault: false},
	}

	for _, role := range roles {
		var existingRole models.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			// Role doesn't exist, create it
			if err := db.Create(&role).Error; err != nil {
				utils.LogError("Failed to create role: " + role.Name)
				return err
			}
		}
	}

	// Seed Claims
	claims := []models.Claim{
		// User claims
		{Name: "user:read"},
		{Name: "user:write"},
		{Name: "user:delete"},
		{Name: "user:update"},
		{Name: "user:refresh"},
		{Name: "user:logout"},

		// Claim claims
		{Name: "claim:read"},
		{Name: "claim:write"},
		{Name: "claim:delete"},
		{Name: "claim:update"},

		// Role claims
		{Name: "role:read"},
		{Name: "role:write"},
		{Name: "role:delete"},
		{Name: "role:update"},
	}

	for _, claim := range claims {
		var existingClaim models.Claim
		if err := db.Where("name = ?", claim.Name).First(&existingClaim).Error; err != nil {
			// Claim doesn't exist, create it
			if err := db.Create(&claim).Error; err != nil {
				utils.LogError("Failed to create claim: " + claim.Name)
				return err
			}
		}
	}

	// Assign all claims to admin role
	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err == nil {
		var allClaims []models.Claim
		db.Find(&allClaims)
		if err := db.Model(&adminRole).Association("Claims").Replace(&allClaims); err != nil {
			utils.LogError("Failed to assign claims to admin role")
			return err
		}
	}

	// Assign basic read claims to user role
	var userRole models.Role
	if err := db.Where("name = ?", "user").First(&userRole).Error; err == nil {
		var userClaims []models.Claim
		db.Where("name IN ?", []string{"user:read", "claim:read", "role:read"}).Find(&userClaims)
		if err := db.Model(&userRole).Association("Claims").Replace(&userClaims); err != nil {
			utils.LogError("Failed to assign claims to user role")
			return err
		}
	}

	// Create default admin user
	var existingAdmin models.User
	if err := db.Where("username = ?", "admin").First(&existingAdmin).Error; err != nil {
		// Admin user doesn't exist, create it
		var adminRoleForUser models.Role
		if err := db.Where("name = ?", "admin").First(&adminRoleForUser).Error; err == nil {
			adminUser := models.User{
				Username: "admin",
				Email:    "admin@knowstack.com",
				Password: "admin123", // This will be hashed by BeforeCreate hook
				Provider: "local",
				RoleID:   adminRoleForUser.ID,
			}
			if err := db.Create(&adminUser).Error; err != nil {
				utils.LogError("Failed to create admin user")
				return err
			}
		}
	}

	utils.LogInfo("Seed data completed successfully")
	return nil
}
