package handler

import (
	"bytes"
	"encoding/json"
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

	db := store.NewStore("./products_copy.json")
	repo := product.NewRepository(db)
	service := product.NewService(repo)
	productHandler := NewProductHandler(service)
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	pr := r.Group("/products")
	{
		pr.GET("/", productHandler.GetAll())
		pr.GET(":id", productHandler.GetByID())
		pr.GET("/search", productHandler.Search())
		pr.POST("/", productHandler.AddProduct())
		pr.DELETE(":id", productHandler.Delete())
		pr.PATCH(":id", productHandler.Patch())
		pr.PUT(":id", productHandler.Put())
	}
	return r
}

func createRequestTest(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", "secret_321")

	return req, httptest.NewRecorder()
}

func TestGetAllProduct_OK(t *testing.T) {
	type ObjTestResponse struct {
		Data []any
	}
	var respTest ObjTestResponse

	r := createServer()

	req, res := createRequestTest(http.MethodGet, "/products/", "")

	r.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)
	err := json.Unmarshal(res.Body.Bytes(), &respTest)
	assert.Nil(t, err)
	assert.True(t, len(respTest.Data) > 0)
}

func TestGetProductById_OK(t *testing.T) {
	type ObjTestResponse struct {
		Data any
	}
	var respTest ObjTestResponse
	r := createServer()

	req, res := createRequestTest(http.MethodGet, "/products/1", "")

	r.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)
	err := json.Unmarshal(res.Body.Bytes(), &respTest)
	assert.Nil(t, err)
	assert.False(t, respTest == ObjTestResponse{})
}

func TestProductHandler_AddProduct(t *testing.T) {
	data := `{"name": "TestPost321", "quantity": 155, "price": 555.99, "code_value": "TFGH312", "expiration": "11/12/1999", "is_published": true }`

	r := createServer()

	req, res := createRequestTest(http.MethodPost, "/products/", data)

	r.ServeHTTP(res, req)

	assert.Equal(t, 201, res.Code)
}

func TestProductHandler_Patch(t *testing.T) {
	data := `{"name": "nombre actualizado 4"}`

	r := createServer()

	req, res := createRequestTest(http.MethodPatch, "/products/502", data)

	r.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)
}

func TestProductHandler_Put(t *testing.T) {
	data := `{"name": "ACTUALIZO NOMBRE", "quantity": 555, "price": 555.99, "code_value": "TFF4455", "expiration": "15/05/2015", "is_published": true }`

	r := createServer()

	req, res := createRequestTest(http.MethodPatch, "/products/504", data)

	r.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)
}

func TestProductHandler_Delete(t *testing.T) {

	r := createServer()

	req, res := createRequestTest(http.MethodDelete, "/products/503", "")

	r.ServeHTTP(res, req)

	assert.Equal(t, 204, res.Code)
}
