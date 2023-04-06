package main

import (
	"github.com/fgiudicatti-meli/web-server/cmd/server/handler"
	"github.com/fgiudicatti-meli/web-server/cmd/server/middlewares"
	"github.com/fgiudicatti-meli/web-server/internal/product"
	"github.com/fgiudicatti-meli/web-server/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	storage := store.NewStore("../../products.json")

	repo := product.NewRepository(storage)
	service := product.NewService(repo)
	productHandler := handler.NewProductHandler(service)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
	products := r.Group("/products")
	{
		products.GET("", productHandler.GetAll())
		products.GET(":id", productHandler.GetByID())
		// --- from here and below middleware apply in every route
		products.Use(middlewares.MiddlewareVerifyToken())
		products.GET("/search", productHandler.Search())
		products.GET("/consumer_price", productHandler.GetPriceProducts())
		products.POST("", productHandler.AddProduct())
		products.DELETE(":id", productHandler.Delete())
		products.PATCH(":id", productHandler.Patch())
		products.PUT(":id", productHandler.Put())
	}

	if err := r.Run(); err != nil {
		panic(err)
	}
}
