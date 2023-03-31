package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
)

type SaludoRequestBody struct {
	Nombre   string `json:"nombre" binding:"required"`
	Apellido string `json:"apellido" binding:"required"`
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

		//bindjson is work as well
		err := json.NewDecoder(c.Request.Body).Decode(&requestBody)
		if err != nil {
			log.Fatal(err)
		}

		//other solution using shouldbindjson or bindjson as well
		/*
			if err := c.ShouldBindJSON(&requestBody); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
			}
		*/
		c.String(200, "Hola %s %s", requestBody.Nombre, requestBody.Apellido)
	})

	router.Run(":8081")
}
