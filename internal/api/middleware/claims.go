package middleware

import (
	"knowstack/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireClaims ensures the authenticated user has ALL of the required claim IDs.
// It assumes JWTMiddleware has already set "claims" in the context.
func RequireClaims(required ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		tokenClaims, ok := raw.(*utils.TokenClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Build a set for quick lookup
		claimSet := make(map[string]struct{}, len(tokenClaims.Claims))
		for _, name := range tokenClaims.Claims {
			claimSet[name] = struct{}{}
		}

		// Verify all required claims are present
		for _, need := range required {
			if _, ok := claimSet[need]; !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
				return
			}
		}

		c.Next()
	}
}
