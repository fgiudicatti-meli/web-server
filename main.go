package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
)

type SaludoRequestBody struct {
	Nombre   string
	Apellido string
}

func main() {
	router := gin.Default()

	router.GET("/hello-world", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.POST("/saludo", func(c *gin.Context) {
		var requestBody SaludoRequestBody
		body, _ := io.ReadAll(c.Request.Body)

		if err := json.Unmarshal(body, &requestBody); err != nil {
			c.JSON(403, gin.H{
				"message": "el cuerpo debe contener las propiedades nombre y apellido en formato de cadena",
			})
		}

		c.String(200, "Hola %s %s", requestBody.Nombre, requestBody.Apellido)
	})

	router.Run(":8081")
}
