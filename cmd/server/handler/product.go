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

type Request struct {
	Name        string  `json:"name,omitempty"`
	Quantity    int     `json:"quantity,omitempty"`
	CodeValue   string  `json:"code_value,omitempty"`
	IsPublished bool    `json:"is_published,omitempty"`
	Expiration  string  `json:"expiration,omitempty"`
	Price       float64 `json:"price,omitempty"`
}

// GetAll documentation with Swagger
// ListProducts godoc
// @Summary List products
// @Tags Products
// @Description get products
// @Produce json
// @Param token header string true "token"
// @Success 200 {object} web.Response
// @Router /products [get]
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

// GetByID documentation with Swagger
// GetByID godoc
// @Summary Get one product
// @Tags Products
// @Description search one product that matches with id
// @Produce json
// @Param id path int true "Product ID"
// @Param token header string true "token"
// @Success 200 {object} web.Response
// @Failure 404 {object} web.ErrorResponse
// @Failure 400 {object} web.ErrorResponse
// @Router /products/{id} [get]
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

// Search documentation with Swagger
// Search godoc
// @Summary search products by price limit
// @Tags Products
// @Description find products that price is bigger than param
// @Produce json
// @Param priceGt query int true "Price"
// @Param token header string true "token"
// @Success 200 {object} web.Response
// @Failure 404 {object} web.ErrorResponse
// @Failure 400 {object} web.ErrorResponse
// @Router /products/search [get]
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

// AddProduct documentation swagger
// AddProduct godoc
// @Summary build a new product
// @Tags Products
// @Description Create a new product and saved in db
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param newBody body domain.Product true "Product"
// @Success 201 {object} web.Response
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Router /products/new [post]
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

// Delete documentation swagger
// Delete godoc
// @Summary eliminate a product
// @Tags Products
// @Description Delete definitive a product
// @Produce json
// @Param token header string true "Token"
// @Param id path int true "ProductID"
// @Success 204 {object} web.Response
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Router /products/{id} [delete]
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

// Put documentation swagger
// Put godoc
// @Summary modify totally a product
// @Tags Products
// @Description update with all fields a product
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param id path int true "Product ID"
// @Param putProduct body domain.Product true "UpdateProduct"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Router /products/{id} [put]
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

// Patch documentation swagger
// Patch godoc
// @Summary Partially update a product
// @Tags Products
// @Description update not totally fields only some
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param id path int true "Product ID"
// @Param patchBody body Request true "updateProduct"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Router /products/{id} [patch]
func (h *productHandler) Patch() gin.HandlerFunc {

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

// PartialUpdate godoc
// @Summary Partially update a product
// @Tags Products
// @Description Update some product fields data
// @Accept json
// @Produce json
// @Param token header string true "Token"
// @Param list query int true "Price"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Router /products/consumer_price [get]

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
