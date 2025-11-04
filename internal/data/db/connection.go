package db

import (
	"errors"
	"knowstack/internal/core/config"
	"knowstack/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Connect(cfg config.Database) error {
	connStr := cfg.ConnectionString()

	var err error
	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return err
	}

	utils.LogInfo("Connected to the database")
	return nil
}

func GetDB() *gorm.DB {
	return db
}

func Close() error {
	sqlDB, err := db.DB()
	if err != nil {
		return errors.New("failed to get the database instance")
	}

	return sqlDB.Close()
}
