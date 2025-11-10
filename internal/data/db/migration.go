package db

import (
	"errors"
	"knowstack/internal/data/models"
	"knowstack/internal/utils"
)

func AutoMigrate() error {
	err := db.AutoMigrate(&models.Role{}, &models.Claim{}, &models.User{})

	if err != nil {
		return errors.New("failed to auto migrate the database")
	}

	utils.LogInfo("Auto migration completed")

	// Run seed data
	if err := SeedData(); err != nil {
		utils.LogError("Failed to seed data")
		return err
	}

	return nil
}
