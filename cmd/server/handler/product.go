package handler

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/fgiudicatti-meli/web-server/internal/domain"
	"github.com/fgiudicatti-meli/web-server/internal/product"
	"github.com/fgiudicatti-meli/web-server/pkg/web"
	"github.com/gin-gonic/gin"
)

type productHandler struct {
	s product.Service
}

// NewProductHandler crea un nuevo controller de productos
func NewProductHandler(s product.Service) *productHandler {
	return &productHandler{
		s: s,
	}
}

// GetAll obtiene todos los productos
func (h *productHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, _ := h.s.GetAll()
		web.Success(c, 200, products)
	}
}

// GetByID obtiene un producto por su id
func (h *productHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}
		productFounded, err := h.s.GetByID(id)
		if err != nil {
			web.Failure(c, 404, errors.New("product not found"))
			return
		}
		web.Success(c, 200, productFounded)
	}
}

// Search busca un producto por precio mayor a un valor
func (h *productHandler) Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		priceParam := c.Query("priceGt")
		price, err := strconv.ParseFloat(priceParam, 64)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid price"))
			return
		}
		products, err := h.s.SearchPriceGt(price)
		if err != nil {
			web.Failure(c, 404, errors.New("product not found"))
			return
		}
		web.Success(c, 200, products)
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
	list := []int{}
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

// Post crear un producto nuevo
func (h *productHandler) AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newProduct domain.Product
		token := c.GetHeader("TOKEN")
		if token == "" {
			web.Failure(c, 401, errors.New("token not found"))
			return
		}
		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("invalid token"))
			return
		}
		err := c.ShouldBindJSON(&newProduct)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid json"))
			return
		}
		valid, err := validateEmptys(&newProduct)
		if !valid {
			web.Failure(c, 400, err)
			return
		}
		valid, err = validateExpiration(newProduct.Expiration)
		if !valid {
			web.Failure(c, 400, err)
			return
		}
		createProduct, err := h.s.Create(newProduct)
		if err != nil {
			web.Failure(c, 400, err)
			return
		}
		web.Success(c, 201, createProduct)
	}
}

// Delete elimina un producto
func (h *productHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("TOKEN")
		if token == "" {
			web.Failure(c, 401, errors.New("token not found"))
			return
		}
		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("invalid token"))
			return
		}
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}
		err = h.s.Delete(id)
		if err != nil {
			web.Failure(c, 404, err)
			return
		}
		web.Success(c, 204, nil)
	}
}

// Put actualiza un producto
func (h *productHandler) Put() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("TOKEN")
		if token == "" {
			web.Failure(c, 401, errors.New("token not found"))
			return
		}
		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("invalid token"))
			return
		}
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}
		_, err = h.s.GetByID(id)
		if err != nil {
			web.Failure(c, 404, errors.New("product not found"))
			return
		}
		if err != nil {
			web.Failure(c, 409, err)
			return
		}
		var productToUpdate domain.Product
		err = c.ShouldBindJSON(&productToUpdate)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid json"))
			return
		}
		valid, err := validateEmptys(&productToUpdate)
		if !valid {
			web.Failure(c, 400, err)
			return
		}
		valid, err = validateExpiration(productToUpdate.Expiration)
		if !valid {
			web.Failure(c, 400, err)
			return
		}
		updateProduct, err := h.s.Update(id, productToUpdate)
		if err != nil {
			web.Failure(c, 409, err)
			return
		}
		web.Success(c, 200, updateProduct)
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
	return func(c *gin.Context) {
		token := c.GetHeader("TOKEN")
		if token == "" {
			web.Failure(c, 401, errors.New("token not found"))
			return
		}
		if token != os.Getenv("TOKEN") {
			web.Failure(c, 401, errors.New("invalid token"))
			return
		}
		var r Request
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}
		_, err = h.s.GetByID(id)
		if err != nil {
			web.Failure(c, 404, errors.New("product not found"))
			return
		}
		if err := c.ShouldBindJSON(&r); err != nil {
			web.Failure(c, 400, errors.New("invalid json"))
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
				web.Failure(c, 400, err)
				return
			}
		}
		p, err := h.s.Update(id, update)
		if err != nil {
			web.Failure(c, 409, err)
			return
		}
		web.Success(c, 200, p)
	}
}
