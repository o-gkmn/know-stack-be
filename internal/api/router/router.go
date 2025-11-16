package router

import (
	"knowstack/internal/api/handlers"
	"knowstack/internal/api/middleware"
	"knowstack/internal/core/services"
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
func NewRouter(service *services.Service) *Router {
	return &Router{
		Handlers: handlers.NewHandlers(service),
		Gin:      gin.New(),
	}
}

func (r *Router) Setup() error {

	// Add CORS middleware
	r.Gin.Use(middleware.CORSMiddleware())

	// Add custom logger middleware
	r.Gin.Use(middleware.LoggerMiddleware())

	// Add custom recovery middleware
	r.Gin.Use(middleware.RecoveryMiddleware())

	utils.LogInfo("Middlewares initialized")

	// Setup the API version 1 routes
	v1 := r.Gin.Group("/api/v1")

	// Setup the routes for the API version 1
	r.setupHealthRoutes(v1)
	r.setupUserRoutes(v1)
	r.setupOAuthRoutes(v1)

	// Setup the swagger routes
	r.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Gin.Static("/docs", "./docs")

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

/*
Setup the user routes for the API version 1
*/
func (r *Router) setupUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/users")
	user.POST("/login", r.Handlers.UserHandler.Login)
	user.POST("/register", r.Handlers.UserHandler.CreateUser)
	user.POST("/refresh", r.Handlers.UserHandler.Refresh)
	user.POST("/logout", r.Handlers.UserHandler.Logout)
	user.POST("/request-password-reset", r.Handlers.UserHandler.RequestPasswordReset)
}

/*
Setup the oauth routes for the API version 1
*/
func (r *Router) setupOAuthRoutes(rg *gin.RouterGroup) {
	oauth := rg.Group("/oauth")
	oauth.GET("/google/login", r.Handlers.OAuthHandler.GoogleLogin)
	oauth.GET("/google/callback", r.Handlers.OAuthHandler.GoogleCallback)
}
