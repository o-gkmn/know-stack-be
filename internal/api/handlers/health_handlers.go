package handlers

import "github.com/gin-gonic/gin"

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// @Summary Check the liveness of the service
// @Description Checks if the service is alive
// @Tags API Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *HealthHandler) CheckLiveness(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "UP", "service": "liveness"})
}
