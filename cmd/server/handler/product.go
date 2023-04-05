package handler

import (
	"encoding/json"
	"fmt"
	"github.com/fgiudicatti-meli/web-server/internal/domain"
	"github.com/fgiudicatti-meli/web-server/internal/product"
	"github.com/fgiudicatti-meli/web-server/pkg/web"
	"net/http"
	"os"
	"strconv"
	"strings"

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
			ctx.JSON(http.StatusUnauthorized, web.NewResponse(http.StatusUnauthorized, nil, "token invalido"))
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

func (c *Product) GetSumProducts() gin.HandlerFunc {
	type response struct {
		Products   any     `json:"products"`
		TotalPrice float64 `json:"total_price"`
	}
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, web.NewResponse(http.StatusUnauthorized, nil, "token invalido"))
			return
		}

		query := ctx.Query("list")
		if query == "" {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "debe ingresar una lista de ids de productos"))
			return
		}
		var filterProducts []domain.Product
		var totalPrice float64
		ids := strings.Split(query, ",")
		for _, idInStr := range ids {
			value, err := strconv.Atoi(idInStr)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "ingrese una lista correct de ids"))
				return
			}
			prd, err := c.service.GetById(value)
			if err != nil {
				ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, nil, "no todos los ids corresponden a un producto"))
				return
			}
			if checkValidate(prd.Id, filterProducts) {
				ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "se encontraron ids repetidos"))
				return
			}
			filterProducts = append(filterProducts, prd)
		}

		for i := range filterProducts {
			totalPrice += filterProducts[i].Price
		}

		switch {
		case len(filterProducts) < 10:
			totalPrice = totalPrice * 1.21
		case len(filterProducts) > 10 && len(filterProducts) < 20:
			totalPrice = totalPrice * 1.17
		default:
			totalPrice = totalPrice * 1.15
		}
		resp := response{Products: filterProducts, TotalPrice: float64(int(totalPrice*100)) / 100}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, resp, ""))
	}
}

func checkValidate(id int, slice []domain.Product) bool {
	for i := range slice {
		if slice[i].Id == id {
			return true
		}
	}

	return false
}
