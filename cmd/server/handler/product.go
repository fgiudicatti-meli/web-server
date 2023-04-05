package handler

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fgiudicatti-meli/web-server/internal/domain"
	"github.com/fgiudicatti-meli/web-server/internal/product"
	"github.com/fgiudicatti-meli/web-server/pkg/web"
	"github.com/gin-gonic/gin"
)

type productHandler struct {
	service product.Service
}

// NewProductHandler crea un nuevo controller de productos
func NewProductHandler(s product.Service) *productHandler {
	return &productHandler{
		service: s,
	}
}

// GetAll obtiene todos los productos
func (h *productHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			web.Failure(ctx, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}
		products, err := h.service.GetAll()
		if err != nil {
			web.Failure(ctx, http.StatusNotFound, errors.New("not found products"))
			return
		}
		web.Success(ctx, 200, products)
	}
}

// GetByID obtiene un producto por su id
func (h *productHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			web.Failure(ctx, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, errors.New("invalid id"))
			return
		}

		productFounded, err := h.service.GetByID(id)
		if err != nil {
			web.Failure(ctx, http.StatusNotFound, errors.New("product not found"))
			return
		}
		web.Success(ctx, http.StatusOK, productFounded)
	}
}

// Search busca un producto por precio mayor a un valor
func (h *productHandler) Search() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			web.Failure(ctx, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, errors.New("invalid price"))
			return
		}

		products, err := h.service.SearchPriceGt(price)
		if err != nil {
			web.Failure(ctx, http.StatusNotFound, errors.New("product not found"))
			return
		}

		web.Success(ctx, http.StatusOK, products)
	}
}

// validateEmptys valida que los campos no esten vacios
func validateEmptys(product *domain.Product) (bool, error) {
	switch {
	case product.Name == "" || product.CodeValue == "" || product.Expiration == "":
		return false, errors.New("fields can't be empty")
	case product.Quantity <= 0 || product.Price <= 0:
		if product.Quantity <= 0 {
			return false, errors.New("quantity must be greater than 0")
		}
		if product.Price <= 0 {
			return false, errors.New("price must be greater than 0")
		}
	}
	return true, nil
}

// validateExpiration valida que la fecha de expiracion sea valida
func validateExpiration(exp string) (bool, error) {
	dates := strings.Split(exp, "/")
	//list := []int{}
	var list []int
	if len(dates) != 3 {
		return false, errors.New("invalid expiration date, must be in format: dd/mm/yyyy")
	}
	for value := range dates {
		number, err := strconv.Atoi(dates[value])
		if err != nil {
			return false, errors.New("invalid expiration date, must be numbers")
		}
		list = append(list, number)
	}
	condition := (list[0] < 1 || list[0] > 31) && (list[1] < 1 || list[1] > 12) && (list[2] < 1 || list[2] > 9999)
	if condition {
		return false, errors.New("invalid expiration date, date must be between 1 and 31/12/9999")
	}
	return true, nil
}

// AddProduct crear un producto nuevo
func (h *productHandler) AddProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			web.Failure(ctx, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		var newProduct domain.Product
		err := ctx.ShouldBindJSON(&newProduct)
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, errors.New("invalid json"))
			return
		}
		valid, err := validateEmptys(&newProduct)
		if !valid {
			web.Failure(ctx, http.StatusBadRequest, err)
			return
		}
		valid, err = validateExpiration(newProduct.Expiration)
		if !valid {
			web.Failure(ctx, http.StatusBadRequest, err)
			return
		}
		createProduct, err := h.service.Create(newProduct)
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, err)
			return
		}

		web.Success(ctx, http.StatusCreated, createProduct)
	}
}

// Delete elimina un producto
func (h *productHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			web.Failure(ctx, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, errors.New("invalid id"))
			return
		}

		err = h.service.Delete(id)
		if err != nil {
			web.Failure(ctx, http.StatusNotFound, err)
			return
		}

		web.Success(ctx, http.StatusNoContent, nil)
	}
}

// Put actualiza un producto
func (h *productHandler) Put() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			web.Failure(ctx, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, errors.New("invalid id"))
			return
		}
		_, err = h.service.GetByID(id)
		if err != nil {
			web.Failure(ctx, http.StatusNotFound, errors.New("product not found"))
			return
		}

		var productToUpdate domain.Product
		err = ctx.ShouldBindJSON(&productToUpdate)
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, errors.New("invalid json"))
			return
		}

		valid, err := validateEmptys(&productToUpdate)
		if !valid {
			web.Failure(ctx, http.StatusBadRequest, err)
			return
		}
		valid, err = validateExpiration(productToUpdate.Expiration)
		if !valid {
			web.Failure(ctx, http.StatusBadRequest, err)
			return
		}

		updateProduct, err := h.service.Update(id, productToUpdate)
		if err != nil {
			web.Failure(ctx, http.StatusConflict, err)
			return
		}

		web.Success(ctx, http.StatusOK, updateProduct)
	}
}

// Patch update selected fields of a product WIP
func (h *productHandler) Patch() gin.HandlerFunc {
	type Request struct {
		Name        string  `json:"name,omitempty"`
		Quantity    int     `json:"quantity,omitempty"`
		CodeValue   string  `json:"code_value,omitempty"`
		IsPublished bool    `json:"is_published,omitempty"`
		Expiration  string  `json:"expiration,omitempty"`
		Price       float64 `json:"price,omitempty"`
	}
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			web.Failure(ctx, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		var r Request
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Failure(ctx, http.StatusBadRequest, errors.New("invalid id"))
			return
		}
		/* other way
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
		*/
		_, err = h.service.GetByID(id)
		if err != nil {
			web.Failure(ctx, http.StatusNotFound, errors.New("product not found"))
			return
		}
		if err := ctx.ShouldBindJSON(&r); err != nil {
			web.Failure(ctx, http.StatusBadRequest, errors.New("invalid json"))
			return
		}

		update := domain.Product{
			Name:        r.Name,
			Quantity:    r.Quantity,
			CodeValue:   r.CodeValue,
			IsPublished: r.IsPublished,
			Expiration:  r.Expiration,
			Price:       r.Price,
		}
		if update.Expiration != "" {
			valid, err := validateExpiration(update.Expiration)
			if !valid {
				web.Failure(ctx, http.StatusBadRequest, err)
				return
			}
		}

		p, err := h.service.Update(id, update)
		if err != nil {
			web.Failure(ctx, http.StatusConflict, err)
			return
		}

		web.Success(ctx, http.StatusOK, p)
	}
}

func (h *productHandler) GetPriceProducts() gin.HandlerFunc {
	type response struct {
		Products   any     `json:"products"`
		TotalPrice float64 `json:"total_price"`
	}
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			web.Failure(ctx, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		query := ctx.Query("list")
		if query == "" {
			web.Failure(ctx, http.StatusBadRequest, errors.New("invalid query param"))
			return
		}
		var filterProducts []domain.Product
		var totalPrice float64
		ids := strings.Split(query, ",")
		for _, idInStr := range ids {
			id, err := strconv.Atoi(idInStr)
			if err != nil {
				web.Failure(ctx, http.StatusBadRequest, errors.New("list of ids invalid"))
				return
			}
			prd, err := h.service.GetByID(id)
			if err != nil {
				web.Failure(ctx, http.StatusBadRequest, errors.New("some ids are not associate with a product"))
				return
			}
			if checkValidate(prd.Id, filterProducts) {
				web.Failure(ctx, http.StatusBadRequest, errors.New("some ids are repeated"))
				return
			}
			if !prd.IsPublished {
				web.Failure(ctx, http.StatusBadRequest, errors.New("remember ids must be a product published"))
				return
			}
			filterProducts = append(filterProducts, prd)
		}

		allRecords, _ := h.service.GetAll()
		if len(allRecords) < len(filterProducts) {
			web.Failure(ctx, http.StatusBadRequest, errors.New("list is too much longer"))
			return
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

		web.Success(ctx, http.StatusOK, resp)
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
