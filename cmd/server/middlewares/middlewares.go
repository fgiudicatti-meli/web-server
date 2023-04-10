package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
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

func CatchPanic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				now := time.Now()
				fmt.Printf("FROM URL: %s\n", ctx.Request.URL.Path)
				fmt.Printf("Verb http: %s\n", ctx.Request.Method)
				fmt.Printf("Weight in Bytes: %b\n", ctx.Request.ContentLength)
				fmt.Printf("Occurs: %s\n", now.Format("2006-01-02 15:04:05"))
			}
		}()

		ctx.Next()
	}
}
