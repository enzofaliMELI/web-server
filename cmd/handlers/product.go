package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/enzofaliMELI/web-server/internal/product/models"
	service "github.com/enzofaliMELI/web-server/internal/product/service"
	"github.com/enzofaliMELI/web-server/pkg/response.go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// -------------------------------- GET Methods --------------------------------

func Pong(ctx *gin.Context) {
	// response
	ctx.String(http.StatusOK, "pong")
}

func GETProducts(ctx *gin.Context) {
	// process
	products := service.GetProducts()

	// response
	ctx.JSON(http.StatusOK, gin.H{"message": "succeed to get all products", "data": products})
}

func ProductsId(ctx *gin.Context) {
	// request
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to parse number", "data": nil})
		return
	}

	// process
	products := service.GetProductById(id)

	// response
	ctx.JSON(http.StatusOK, gin.H{"message": "succeed to get all products", "data": products})
}

func ProductsSearch(ctx *gin.Context) {
	// request
	price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to parse number", "data": nil})
		return
	}

	// process
	products := service.GetProductByPriceGt(price)

	// response
	ctx.JSON(http.StatusOK, gin.H{"message": "succeed to get all products", "data": products})
}

// -------------------------------- POST Methods --------------------------------

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

func POSTProducts(ctx *gin.Context) {
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

	// The services recieve a product struct
	product := models.Product{
		Id:           0,
		Name:         request.Name,
		Quantity:     request.Quantity,
		Code_value:   request.Code_value,
		Is_published: request.Is_published,
		Expiration:   request.Expiration,
		Price:        request.Price,
	}

	if service.InvalidCodeValue(product, service.Products) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "there is already a product with that code", "data": nil})
		return
	}

	if service.InvalidExpiration(product) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid expiration date, (format: DD/MM/YYYY)", "data": nil})
		return
	}

	// process
	service.LastID++
	product.Id = service.LastID
	service.Products = append(service.Products, product)

	// response
	ctx.JSON(http.StatusCreated, response.Ok("succeed to upload a product", request))
}
