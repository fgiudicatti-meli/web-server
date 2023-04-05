package test

import (
	"bytes"
	"encoding/json"
	"github.com/fgiudicatti-meli/web-server/cmd/server/handler"
	"github.com/fgiudicatti-meli/web-server/internal/product"
	"github.com/fgiudicatti-meli/web-server/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func createServer() *gin.Engine {
	_ = os.Setenv("TOKEN", "secret_321")
	db := store.NewStore("../products.json")
	repo := product.NewRepository(db)
	service := product.NewService(repo)
	p := handler.NewProduct(service)
	r := gin.Default()

	pr := r.Group("/products")
	{
		pr.GET("/", p.GetAll())
		pr.GET("/:id", p.GetById())
		pr.POST("/", p.Save())
		pr.DELETE("/:id", p.Delete())
	}

	return r
}

func createRequest(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", "secret_321")

	return req, httptest.NewRecorder()
}

func TestGetAllProductsOK(t *testing.T) {
	type TestRes struct {
		Data []any
	}
	var resp TestRes
	r := createServer()

	req, res := createRequest(http.MethodGet, "/products/", "")

	r.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)
	err := json.Unmarshal(res.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.True(t, len(resp.Data) == 500)
}

func TestGetProductByIdOK(t *testing.T) {
	r := createServer()

	req, res := createRequest(http.MethodGet, "/products/100", "")

	r.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)
}

func TestAddProductOK(t *testing.T) {
	r := createServer()

	req, res := createRequest(http.MethodPost, "/products/", `{
"name": "Francisco test 2", "quantity": 500, "price": 256.85, "is_published": true, "expiration": "15/09/2022", "code_value": "F1234"}`)

	r.ServeHTTP(res, req)

	assert.Equal(t, 201, res.Code)
}

func TestDeleteProductOK(t *testing.T) {
	r := createServer()

	req, res := createRequest(http.MethodDelete, "/products/100", "")

	r.ServeHTTP(res, req)

	assert.Equal(t, 204, res.Code)
}


