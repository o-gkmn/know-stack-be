package api

import (
	"errors"
	"fmt"
	"knowstack/internal/api/router"
	"knowstack/internal/core/config"
	"knowstack/internal/core/services"
	"knowstack/internal/data/db"
	"knowstack/internal/utils"

	"gorm.io/gorm"
)

type Server struct {
	Config  config.Server
	Router  *router.Router
	DB      *gorm.DB
}

/*
Create a new server instance
Returns a pointer to the server instance
*/
func NewServer(config config.Server) *Server {
	//Connect to the database
	err := db.Connect(config.Database)
	if err != nil {
		utils.LogFatalWithErr("Failed to connect to the database", err)
	}

	// Auto migrate the database
	err = db.AutoMigrate()
	if err != nil {
		utils.LogFatalWithErr("Failed to auto migrate the database", err)
	}

	// Create a new service instance
	serviceInstance := services.NewService(db.GetDB())

	// Create a new router instance and setup the routes
	r := router.NewRouter(serviceInstance)
	r.Setup()

	utils.LogInfo("Server initialized")

	return &Server{
		Config:  config,
		Router:  r,
		DB:      db.GetDB(),
	}
}

/*
Check if the server is ready to start
Returns true if the server is ready, false otherwise
*/
func (s *Server) Ready() bool {
	return s.Router != nil && s.DB != nil
}

/*
Start the server on the configured port and host
Returns an error if the server is not ready
*/
func (s *Server) Start() error {
	if !s.Ready() {
		err := errors.New("server is not ready")
		utils.LogError("Server is not ready to start")
		return err
	}

	addr := fmt.Sprintf("%s:%s", s.Config.Host, s.Config.Port)
	utils.LogInfo("Starting server", "address", addr)

	return s.Router.Gin.Run(addr)
}
