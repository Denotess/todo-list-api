package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"main.go/internal/helpers"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		claims, err := helpers.VerifyToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.Subject)
		ctx.Next()
	}
}
