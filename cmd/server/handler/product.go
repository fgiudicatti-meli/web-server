package handler

import (
	"encoding/json"
	"fmt"
	"github.com/fgiudicatti-meli/web-server/internal/product"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Name        string  `json:"name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	CodeValue   string  `json:"code_value" binding:"required"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}

type Product struct {
	service product.Service
}

func NewProduct(p product.Service) *Product {
	return &Product{
		service: p,
	}
}

func (c *Product) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(401, gin.H{
				"error": "token inválido",
			})
			return
		}
		allProducts, err := c.service.GetAll()
		if err != nil {
			ctx.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, allProducts)
	}
}

func (c *Product) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(401, gin.H{
				"error": "token inválido",
			})
			return
		}
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "id invalido",
			})
			return
		}
		productById, err := c.service.GetById(id)
		if err != nil {
			ctx.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, productById)
	}
}

func (c *Product) Save() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(401, gin.H{"error": "token inválido"})
			return
		}
		var req Request
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}
		p, err := c.service.Save(req.Name, req.CodeValue, req.Expiration, req.Quantity, req.Price, req.IsPublished)
		if err != nil {
			ctx.JSON(404, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, p)
	}
}

func (c *Product) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(401, gin.H{
				"error": "token invalido",
			})
			return
		}

		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "id invalido",
			})
			return
		}

		var req Request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		switch {
		case req.Name == "":
			ctx.JSON(400, gin.H{
				"error": "nombre es requerido",
			})
			return
		case req.CodeValue == "":
			ctx.JSON(400, gin.H{
				"error": "codigo valor es requerido",
			})
			return
		case req.Expiration == "":
			ctx.JSON(400, gin.H{
				"error": "expiracion es requerido",
			})
			return
		case req.Price == 0:
			ctx.JSON(400, gin.H{
				"error": "precio es requerido",
			})
			return
		case req.Quantity == 0:
			ctx.JSON(400, gin.H{
				"error": "cantidad es requerido",
			})
			return
		}

		updateProduct, err := c.service.Update(int(id), req.Name, req.CodeValue, req.Expiration, req.Quantity, req.Price, req.IsPublished)
		if err != nil {
			ctx.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, updateProduct)
	}
}

func (c *Product) UpdateName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(401, gin.H{
				"error": "token invalido",
			})
			return
		}

		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "id invalido",
			})
			return
		}
		var req Request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		if req.Name == "" {
			ctx.JSON(400, gin.H{
				"error": "nombre requerido",
			})
			return
		}

		productUpdateByName, err := c.service.UpdateName(int(id), req.Name)
		if err != nil {
			ctx.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(200, productUpdateByName)
	}
}

func (c *Product) UpdatePartial() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "token invalido",
			})
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request",
			})
			return
		}

		oldProduct, err := c.service.GetById(id)
		if err != nil {
			ctx.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := json.NewDecoder(ctx.Request.Body).Decode(&oldProduct); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		oldProduct.Id = id

		prd, err := c.service.Update(id, oldProduct.Name, oldProduct.CodeValue, oldProduct.Expiration, oldProduct.Quantity, oldProduct.Price, oldProduct.IsPublished)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data":    prd,
		})

	}
}

func (c *Product) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(401, gin.H{
				"error": "token invalido",
			})
			return
		}

		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "id invalido",
			})
			return
		}

		err = c.service.Delete(int(id))
		if err != nil {
			ctx.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"data": fmt.Sprintf("El producto %d ha sido eliminado", id),
		})
	}
}
