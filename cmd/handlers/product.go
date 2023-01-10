package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/enzofaliMELI/web-server/internal/product"
	"github.com/enzofaliMELI/web-server/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

const secretKey = "1234"

var (
	ErrUnauthorized = errors.New("error: invalid token")
)

type Request struct {
	Name         string  `json:"name" validate:"required"`
	Quantity     int     `json:"quantity" validate:"required"`
	Code_value   string  `json:"code_value" validate:"required"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration" validate:"required"`
	Price        float64 `json:"price" validate:"required"`
}

// -------------------------------- GET Methods --------------------------------

func Pong(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

func GetAllProducts(ctx *gin.Context) {
	products := product.GetProducts()

	ctx.JSON(http.StatusOK, gin.H{"message": "succeed to get all products", "data": products})
}

func GetProductsId(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to parse number", "data": nil})
		return
	}

	products := product.GetProductById(id)

	ctx.JSON(http.StatusOK, gin.H{"message": "succeed to get all products", "data": products})
}

func GetProductsSearch(ctx *gin.Context) {
	price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to parse number", "data": nil})
		return
	}

	products := product.GetProductByPriceGt(price)

	ctx.JSON(http.StatusOK, gin.H{"message": "succeed to get all products", "data": products})
}

// -------------------------------- POST Methods --------------------------------

func StoreProduct(ctx *gin.Context) {
	// request
	var request Request

	token := ctx.GetHeader("token")
	if token != secretKey {
		ctx.JSON(http.StatusUnauthorized, response.Err(ErrUnauthorized))
		return
	}

	err := ctx.ShouldBindJSON(&request)
	fmt.Println(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": nil})
		return
	}

	validate := validator.New()
	if err := validate.Struct(&request); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error(), "data": nil})
		return
	}

	if product.InvalidCodeValue(request.Code_value) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "there is already a product with that code", "data": nil})
		return
	}

	if product.InvalidExpiration(request.Expiration) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid expiration date, (format: DD/MM/YYYY)", "data": nil})
		return
	}

	// process
	prod := product.Product{
		Id:           0,
		Name:         request.Name,
		Quantity:     request.Quantity,
		Code_value:   request.Code_value,
		Is_published: request.Is_published,
		Expiration:   request.Expiration,
		Price:        request.Price,
	}

	product.LastID++
	prod.Id = product.LastID
	product.Products = append(product.Products, prod)

	// response
	ctx.JSON(http.StatusCreated, response.Ok("succeed to upload a product", request))
}
