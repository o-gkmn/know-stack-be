package middleware

import (
	"knowstack/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs HTTP requests using custom slog logger
func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Start timer
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		// Proccess request
		ctx.Next()

		// Calculate latency
		latency := time.Since(start)

		// Client IP
		clientIP := ctx.ClientIP()

		// Method
		method := ctx.Request.Method

		// Status code
		statusCode := ctx.Writer.Status()

		// Error message
		var errorMsg string
		if len(ctx.Errors) > 0 {
			errorMsg = ctx.Errors.String()
		}

		// Build arguments
		logArgs := []interface{}{
			"client_ip", clientIP,
			"latency", latency,
			"method", method,
			"path", path,
			"status", statusCode,
		}

		// Add query string if exists
		if raw != "" {
			logArgs = append(logArgs, "query", raw)
		}

		// Add error if exists
		if errorMsg != "" {
			logArgs = append(logArgs, "error", errorMsg)
		}

		// Log based on status code
		if statusCode >= 500 {
			utils.LogError("HTTP request", logArgs...)
		} else if statusCode >= 400 {
			utils.LogWarn("HTTP request", logArgs...)
		} else {
			utils.LogInfo("HTTP request", logArgs...)
		}
	}
}

// RecoveryMiddleware recovers from panics and logs them using custom slog logger
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recoverd interface{}) {
		utils.LogError("Panic recovered",
			"error", recoverd,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"client_ip", c.ClientIP(),
		)
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Internal server error",
		})
	})
}
