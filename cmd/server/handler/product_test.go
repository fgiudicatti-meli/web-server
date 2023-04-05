package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fgiudicatti-meli/web-server/cmd/server/handler"
	"github.com/fgiudicatti-meli/web-server/internal/domain"
	"github.com/fgiudicatti-meli/web-server/internal/product"
	"github.com/fgiudicatti-meli/web-server/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type response struct {
	Data interface{} `json:"data"`
}

func createServer(token string) *gin.Engine {

	if token != "" {
		err := os.Setenv("TOKEN", token)
		if err != nil {
			panic(err)
		}
	}

	db := store.NewStore("./products_copy.json")
	repo := product.NewRepository(db)
	service := product.NewService(repo)
	productHandler := handler.NewProductHandler(service)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	pr := r.Group("/products")
	{
		pr.GET("", productHandler.GetAll())
		pr.GET(":id", productHandler.GetByID())
		pr.GET("/search", productHandler.Search())
		pr.POST("", productHandler.AddProduct())
		pr.DELETE(":id", productHandler.Delete())
		pr.PATCH(":id", productHandler.Patch())
		pr.PUT(":id", productHandler.Put())
	}
	return r
}

func createRequestTest(method string, url string, body string, token string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	if token != "" {
		req.Header.Add("TOKEN", token)
	}
	return req, httptest.NewRecorder()
}

func loadProducts(path string) ([]domain.Product, error) {
	var products []domain.Product
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func writeProducts(path string, list []domain.Product) error {
	listToBytes, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, listToBytes, 0644)
	if err != nil {
		return err
	}
	return err
}

func Test_GetAll_OK(t *testing.T) {
	var expectd = response{Data: []domain.Product{}}

	r := createServer("my-secret-token")
	req, rr := createRequestTest(http.MethodGet, "/products", "", "my-secret-token")

	p, err := loadProducts("./products_copy.json")
	if err != nil {
		panic(err)
	}
	expectd.Data = p
	actual := map[string][]domain.Product{}

	r.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expectd.Data, actual["data"])
}

func Test_GetOne_OK(t *testing.T) {
	var expectd = response{Data: domain.Product{
		Id:          1,
		Name:        "Oil - Margarine",
		Quantity:    439,
		CodeValue:   "S82254D",
		IsPublished: true,
		Expiration:  "15/12/2021",
		Price:       71.42,
	}}

	r := createServer("my-secret-token")
	req, rr := createRequestTest(http.MethodGet, "/products/1", "", "my-secret-token")
	r.ServeHTTP(rr, req)

	p, err := loadProducts("./products_copy.json")
	if err != nil {
		panic(err)
	}
	expectd.Data = p[0]
	actual := map[string]domain.Product{}

	assert.Equal(t, 200, rr.Code)
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expectd.Data, actual["data"])
}

func Test_Post_OK(t *testing.T) {
	var expectd = response{Data: domain.Product{
		Id:          500,
		Name:        "Oil - Margarine",
		Quantity:    439,
		CodeValue:   "TEST45050",
		IsPublished: true,
		Expiration:  "15/12/2021",
		Price:       50.50,
	}}

	productTest, _ := json.Marshal(expectd.Data)

	r := createServer("my-secret-token")
	req, rr := createRequestTest(http.MethodPost, "/products", string(productTest), "my-secret-token")

	p, _ := loadProducts("./products_copy.json")

	r.ServeHTTP(rr, req)
	actual := map[string]domain.Product{}
	_ = json.Unmarshal(rr.Body.Bytes(), &actual)
	_ = writeProducts("./products_copy.json", p)

	assert.Equal(t, 201, rr.Code)
	assert.Equal(t, expectd.Data, actual["data"])

}

func Test_Delete_OK(t *testing.T) {

	r := createServer("my-secret-token")
	req, rr := createRequestTest(http.MethodDelete, "/products/1", "", "my-secret-token")

	p, err := loadProducts("./products_copy.json")
	if err != nil {
		panic(err)
	}

	r.ServeHTTP(rr, req)

	err = writeProducts("./products_copy.json", p)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, 204, rr.Code)
	assert.Nil(t, rr.Body.Bytes())
}

func Test_BadRequest(t *testing.T) {

	test := []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete}

	r := createServer("my-secret-token")
	for _, method := range test {
		req, rr := createRequestTest(method, "/products/not_number", "", "my-secret-token")
		r.ServeHTTP(rr, req)
		assert.Equal(t, 400, rr.Code)
	}

}

func Test_NotFound(t *testing.T) {

	test := []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete}

	r := createServer("my-secret-token")
	for _, method := range test {
		req, rr := createRequestTest(method, "/products/1000", "{}", "my-secret-token")
		r.ServeHTTP(rr, req)
		assert.Equal(t, 404, rr.Code)
	}
}

func Test_Unauthorized(t *testing.T) {

	test := []string{http.MethodPut, http.MethodPatch, http.MethodDelete}

	r := createServer("my-secret-token")
	for _, method := range test {
		req, rr := createRequestTest(method, "/products/10", "{}", "not-my-token")
		r.ServeHTTP(rr, req)
		assert.Equal(t, 401, rr.Code)
	}
}

func Test_Post_Unauthorized(t *testing.T) {
	r := createServer("my-secret-token")
	req, rr := createRequestTest(http.MethodPost, "/products", "{}", "not-my-token")
	r.ServeHTTP(rr, req)
	assert.Equal(t, 401, rr.Code)
}
