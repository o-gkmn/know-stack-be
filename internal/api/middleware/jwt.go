package middleware

import (
	"knowstack/internal/utils"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		utils.LogInfo("JWT Middleware")

		token := ctx.GetHeader("Authorization")
		utils.LogInfo("Token: %+v", token)

		if token == "" {
			utils.LogInfo("token is empty")
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		token = utils.ExtractBearerToken(token)
		if token == "" {
			utils.LogInfo("token is empty")
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		claims, err := utils.VerifyAccessToken(token)
		if err != nil {
			utils.LogErrorWithErr("Failed to verify JWT", err)
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
