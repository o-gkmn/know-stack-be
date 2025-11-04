package router

import (
	"knowstack/internal/api/handlers"
	"knowstack/internal/api/middleware"
	"knowstack/internal/utils"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	Handlers *handlers.Handlers
	Gin      *gin.Engine
}

/*
Creates a new router instance
*/
func NewRouter() *Router {
	return &Router{
		Handlers: handlers.NewHandlers(),
		Gin:      gin.New(),
	}
}

func (r *Router) Setup() error {
	// Add custom logger middleware
	r.Gin.Use(middleware.LoggerMiddleware())

	// Add custom recovery middleware
	r.Gin.Use(middleware.RecoveryMiddleware())

	utils.LogInfo("Middlewares initialized")

	// Setup the API version 1 routes
	v1 := r.Gin.Group("/api/v1")

	// Setup the routes for the API version 1
	r.setupHealthRoutes(v1)

	// Setup the swagger routes
	r.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	utils.LogInfo("Routes initialized")
	return nil
}

/*
Setup the health routes for the API version 1
*/
func (r *Router) setupHealthRoutes(rg *gin.RouterGroup) {
	health := rg.Group("/health")
	health.GET("", r.Handlers.HealthHandler.CheckLiveness)
}
