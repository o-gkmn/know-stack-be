package db

import (
	"errors"
	"knowstack/internal/utils"
)

func AutoMigrate() error {
	err := db.AutoMigrate()

	if err != nil {
		return errors.New("failed to auto migrate the database")
	}

	utils.LogInfo("Auto migration completed")

	return nil
}
