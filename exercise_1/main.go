package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
)

type Product struct {
	Id          int     `json:"id" binding:"required,unique"`
	Name        string  `json:"name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	CodeValue   string  `json:"code_value" binding:"required"`
	IsPublished bool    `json:"is_published" binding:"required" default:"false"`
	Expiration  string  `json:"expiration" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
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

func GetLastId() int {
	return products[len(products)-1].Id
}

func AddProduct(c *gin.Context) {
	type Request struct {
		Name        string  `json:"name" binding:"required"`
		Quantity    int     `json:"quantity" binding:"required"`
		CodeValue   string  `json:"code_value" binding:"required"`
		IsPublished bool    `json:"is_published" binding:"required"`
		Expiration  string  `json:"expiration" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
	}
	var req Request

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	newProduct := Product{
		Id:          GetLastId() + 1,
		Name:        req.Name,
		Quantity:    req.Quantity,
		CodeValue:   req.CodeValue,
		IsPublished: req.IsPublished,
		Expiration:  req.Expiration,
		Price:       req.Price,
	}

	products = append(products, newProduct)
	c.JSON(201, newProduct)
}

func main() {
	GetAll()
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	GroupProductEndpoints := router.Group("/productos")
	GroupProductEndpoints.GET("/", FilterProductsByPrice)
	GroupProductEndpoints.POST("/", AddProduct)
	GroupProductEndpoints.GET("/:id", FindProductById)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
