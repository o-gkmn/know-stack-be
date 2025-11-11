package main

import (
	"knowstack/internal/api"
	"knowstack/internal/core/config"
	"knowstack/internal/core/logging"
	"knowstack/internal/utils"

	_ "knowstack/docs"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

// @title Knowstack API
// @version 1.0
// @description API for the Knowstack project
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Gin is set to release mode to get rid of the debug logs
	gin.SetMode(gin.ReleaseMode)

	// loading env file from the root folder
	_ = gotenv.Load(".env")

	// initialize config
	cfg := config.DefaultServerConfigFromEnv()

	// initialize logger with config
	logging.Init(cfg.Logger)
	utils.LogInfo("Logger initialized")

	// Initializes the server instance and register the routes
	s := api.NewServer(cfg)

	s.Start()
}
