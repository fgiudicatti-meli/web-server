package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

var products []Product

func GetAll() {
	productsFromJson, err := os.ReadFile("./products.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(productsFromJson, &products); err != nil {
		log.Fatal(err)
	}
}

func FindProductById(c *gin.Context) {
	for i := range products {
		if strconv.Itoa(products[i].Id) == c.Param("id") {
			c.JSON(200, products[i])
			return
		}
	}

	c.JSON(404, gin.H{
		"message": "User not found",
	})
}

func FilterProductsByPrice(c *gin.Context) {
	GetAll()
	var filteredProducts []Product
	formatQueryParam, _ := strconv.ParseFloat(c.Query("priceGt"), 64)

	for i := range products {
		if products[i].Price > formatQueryParam {
			filteredProducts = append(filteredProducts, products[i])
		}
	}

	if len(filteredProducts) == 0 {
		c.JSON(404, gin.H{
			"message": "Results are empty",
		})
	} else {
		c.JSON(200, gin.H{
			"results": len(filteredProducts),
			"data":    filteredProducts,
		})
	}

}

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.GET("productos", FilterProductsByPrice)
	router.GET("productos/:id", FindProductById)
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
