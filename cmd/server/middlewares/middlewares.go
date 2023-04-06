package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// MiddlewareVerifyToken gestiona token entre request
func MiddlewareVerifyToken() gin.HandlerFunc {
	//add token in the header
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
