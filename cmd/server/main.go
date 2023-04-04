package main

import (
	"github.com/fgiudicatti-meli/web-server/cmd/server/handler"
	"github.com/fgiudicatti-meli/web-server/internal/product"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {
	_ = godotenv.Load()
	repo := product.NewRepository()
	service := product.NewService(repo)
	p := handler.NewProduct(service)

	r := gin.Default()
	productsGroup := r.Group("/products")
	{
		productsGroup.POST("/", p.Save())
		productsGroup.GET("/", p.GetAll())
		productsGroup.GET("/:id", p.GetById())
		productsGroup.PUT("/:id", p.Update())
		productsGroup.PATCH("/:id/name", p.UpdateName())
		productsGroup.PATCH("/:id", p.UpdatePartial())
		productsGroup.DELETE("/:id", p.Delete())
	}

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
