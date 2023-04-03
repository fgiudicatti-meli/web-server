package main

import (
	"github.com/fgiudicatti-meli/web-server/package_oriented_design/cmd/server/handler"
	"github.com/fgiudicatti-meli/web-server/package_oriented_design/internal/product"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := product.NewRepository()
	service := product.NewService(repo)
	p := handler.NewProduct(service)

	r := gin.Default()
	productsGroup := r.Group("/products")
	productsGroup.POST("/", p.Save())
	productsGroup.GET("/", p.GetAll())
	productsGroup.GET("/:id", p.GetById())

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
