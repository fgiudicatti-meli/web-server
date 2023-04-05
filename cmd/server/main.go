package main

import (
	"github.com/fgiudicatti-meli/web-server/cmd/server/handler"
	"github.com/fgiudicatti-meli/web-server/internal/product"
	"github.com/fgiudicatti-meli/web-server/pkg/store"
	"github.com/joho/godotenv"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	_ = godotenv.Load()

	db := store.NewStore("products.json")
	if err := db.Check(); err != nil {
		log.Fatal("error al intentar cargar el archivo del store")
	}
	repo := product.NewRepository(db)
	service := product.NewService(repo)
	p := handler.NewProduct(service)

	r := gin.Default()
	productsGroup := r.Group("/products")
	{
		productsGroup.GET("/", p.GetAll())
		productsGroup.GET("/consumer_price", p.GetPriceProducts())
		productsGroup.GET("/:id", p.GetById())
		productsGroup.POST("/", p.Save())
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
