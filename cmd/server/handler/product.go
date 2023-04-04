package handler

import (
	"encoding/json"
	"fmt"
	"github.com/fgiudicatti-meli/web-server/internal/product"
	"github.com/fgiudicatti-meli/web-server/pkg/web"
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
			ctx.JSON(http.StatusUnauthorized, web.NewResponse(http.StatusUnauthorized, nil, "token invalido"))
			return
		}
		allProducts, err := c.service.GetAll()
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, nil, err.Error()))
			return
		}
		ctx.JSON(200, web.NewResponse(http.StatusOK, allProducts, ""))
	}
}

func (c *Product) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, nil, "token invalido"))
			return
		}
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "id invalido"))
			return
		}
		productById, err := c.service.GetById(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, productById, ""))
	}
}

func (c *Product) Save() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, nil, "token invalido"))
			return
		}

		var req Request
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, err.Error()))
			return
		}

		createProduct, err := c.service.Save(req.Name, req.CodeValue, req.Expiration, req.Quantity, req.Price, req.IsPublished)
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, nil, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, createProduct, ""))
	}
}

func (c *Product) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, web.NewResponse(http.StatusUnauthorized, nil, "token invalido"))
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "id invalido"))
			return
		}

		var req Request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, err.Error()))
			return
		}

		switch {
		case req.Name == "":
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "nombre es requerido"))
			return
		case req.CodeValue == "":
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "codeValue es requerido"))
			return
		case req.Expiration == "":
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "expiration es requerido"))
			return
		case req.Price == 0:
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "price es requerido"))
			return
		case req.Quantity == 0:
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "quantity es requerido"))
			return
		}

		updateProduct, err := c.service.Update(id, req.Name, req.CodeValue, req.Expiration, req.Quantity, req.Price, req.IsPublished)
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, updateProduct, ""))
	}
}

func (c *Product) UpdateName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, web.NewResponse(http.StatusUnauthorized, nil, "token invalido"))
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "id invalido"))
			return
		}
		var req Request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, err.Error()))
			return
		}

		if req.Name == "" {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "nombre es requerido"))
			return
		}

		productUpdateByName, err := c.service.UpdateName(id, req.Name)
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, nil, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, productUpdateByName, ""))
	}
}

func (c *Product) UpdatePartial() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, web.NewResponse(http.StatusUnauthorized, nil, "token invalido"))
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "id invalido"))
			return
		}

		oldProduct, err := c.service.GetById(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, nil, err.Error()))
			return
		}

		if err := json.NewDecoder(ctx.Request.Body).Decode(&oldProduct); err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "bad request"))
			return
		}
		oldProduct.Id = id

		prd, err := c.service.Update(id, oldProduct.Name, oldProduct.CodeValue, oldProduct.Expiration, oldProduct.Quantity, oldProduct.Price, oldProduct.IsPublished)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, prd, ""))

	}
}

func (c *Product) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, web.NewResponse(http.StatusUnauthorized, nil, "token invalido"))
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "id invalido"))
			return
		}

		err = c.service.Delete(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, fmt.Sprintf("El producto con id %d ha sido eliminado", id), ""))
	}
}
