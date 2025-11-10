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
			utils.LogInfo("Created role: " + role.Name)
		}
	}

	utils.LogInfo("Seed data completed successfully")
	return nil
}
