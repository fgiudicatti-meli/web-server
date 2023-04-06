package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// MiddlewareVerifyToken check if a token is valid 
func MiddlewareVerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}
		ctx.Next()
	}
}
