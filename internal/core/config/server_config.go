package config

import (
	"knowstack/internal/utils"
)

type Server struct {
	Port     string
	Host     string
	Database Database
	Logger   Logger
}

type Logger struct {
	Level       string
	Format      string
	Environment string
}

/*
Load configurations from the .env file
and use the provided value as a fallback
if settings is not defined in the .env file
*/
func DefaultServerConfigFromEnv() Server {
	return Server{
		Port: utils.GetEnv("PORT", "8080"),
		Host: utils.GetEnv("HOST", "0.0.0.0"),
		Database: Database{
			Host:     utils.GetEnv("DB_HOST", "localhost"),
			Port:     utils.GetEnv("DB_PORT", "5432"),
			User:     utils.GetEnv("DB_USER", "postgres"),
			Password: utils.GetEnv("DB_PASSWORD", "postgres"),
			Database: utils.GetEnv("DB_NAME", "knowstack"),
			SSLMode:  utils.GetEnv("DB_SSLMODE", "disable"),
		},
		Logger: Logger{
			Level:       utils.GetEnv("LOG_LEVEL", "info"),
			Format:      utils.GetEnv("LOG_FORMAT", "text"),
			Environment: utils.GetEnv("LOG_ENVIRONMENT", "development"),
		},
	}
}
