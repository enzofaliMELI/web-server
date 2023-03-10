package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/enzofaliMELI/web-server/cmd/server/middleware"
	"github.com/enzofaliMELI/web-server/internal/domain"
	"github.com/enzofaliMELI/web-server/internal/product"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func createServerProductsTest() *gin.Engine {
	// Read all files
	db := []domain.Product{}
	product.OpenProducts(&db)

	// Server
	server := gin.Default()

	// Router
	repository := product.NewRepository(&db)
	service := product.NewService(repository)
	handler := NewProduct(service)

	products := server.Group("/products")
	{
		products.GET("/", handler.GetAll())
		products.GET("/:id", handler.GetById())
		products.GET("/search", handler.GetPriceGt())
		products.POST("/", middleware.TokenAuthMiddleware(), handler.Store())
		products.PUT("/:id", middleware.TokenAuthMiddleware(), handler.UpdateProduct())
		products.PATCH("/:id", middleware.TokenAuthMiddleware(), handler.UpdatePATCH())
		products.DELETE("/:id", middleware.Middlewares(handler.DeleteProduct())...)
	}

	return server
}

func createRequestTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	request.Header.Add("token", "1234") // No funciona

	return request, httptest.NewRecorder()
}

// Test de Exito ---------------------------------------------------------------------------------
func Test_GetAll(t *testing.T) {
	// Arrange
	server := createServerProductsTest()
	request, response := createRequestTest(http.MethodGet, "/products/", "")

	// Act
	server.ServeHTTP(response, request)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, len(body) > 0)
}

func Test_GetById(t *testing.T) {
	// Arrange
	server := createServerProductsTest()
	request, response := createRequestTest(http.MethodGet, "/products/1", "")

	// Act
	server.ServeHTTP(response, request)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, len(body) > 0)
}

/*
func Test_Store(t *testing.T) {
	// Arrange
	server := createServerProductsTest()
	request, response := createRequestTest(http.MethodPost, "/products/", `{"name":"Chicken - Soup Base","quantity":479,"code_value":"0j3","is_published":false,"expiration":"11/12/2021","price":515.93}`)

	// Act
	server.ServeHTTP(response, request)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, len(body) > 0)
	assert.Equal(t, "1234", response.Header().Get("token"))
}

func Test_Delete(t *testing.T) {
	// Arrange
	server := createServerProductsTest()
	request, response := createRequestTest(http.MethodPost, "/products/10", "")

	// Act
	server.ServeHTTP(response, request)

	// Assert
	assert.Equal(t, http.StatusOK, response.Code)
}

// Test de Fallo ---------------------------------------------------------------------------------
func Test_GetById_BadReq(t *testing.T) {
	// Arrange
	server := createServerProductsTest()
	request, response := createRequestTest(http.MethodGet, "/products/:owi0+", "")

	// Act
	server.ServeHTTP(response, request)

	// Assert
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func Test_Delete_BadReq(t *testing.T) {
	// Arrange
	server := createServerProductsTest()
	request, response := createRequestTest(http.MethodPost, "/products/:owi0+", "")

	// Act
	server.ServeHTTP(response, request)

	// Assert
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func Test_Store_BadReq(t *testing.T) {
	// Arrange
	server := createServerProductsTest()
	request, response := createRequestTest(http.MethodPost, "/products/", `{`)

	// Act
	server.ServeHTTP(response, request)

	// Assert
	assert.Equal(t, http.StatusBadRequest, response.Code)
}
*/
