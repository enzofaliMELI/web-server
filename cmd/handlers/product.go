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

type Product struct {
	s product.Service
}

type Request struct {
	Name         string  `json:"name" validate:"required"`
	Quantity     int     `json:"quantity" validate:"required"`
	Code_value   string  `json:"code_value" validate:"required"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration" validate:"required"`
	Price        float64 `json:"price" validate:"required"`
}

func NewProduct(s product.Service) *Product {
	return &Product{s: s}
}

// -------------------------------- GET Methods --------------------------------

func Pong(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

func (p *Product) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request

		// Process
		products, err := p.s.GetAll()
		if err != nil {
			ctx.JSON(500, nil)
			return
		}

		// Response
		ctx.JSON(http.StatusOK, gin.H{"message": "succeed to get all products", "data": products})
	}
}

func (p *Product) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to parse number", "data": nil})
			return
		}
		// Process
		products, err := p.s.GetById(id)
		if err != nil {
			ctx.JSON(500, nil)
			return
		}

		// Response
		ctx.JSON(http.StatusOK, gin.H{"message": "succeed to get all products", "data": products})
	}
}

func (p *Product) GetPriceGt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to parse number", "data": nil})
			return
		}
		// Process
		products, err := p.s.GetPriceGt(price)
		if err != nil {
			ctx.JSON(500, nil)
			return
		}

		// Response
		ctx.JSON(http.StatusOK, gin.H{"message": "succeed to get all products", "data": products})
	}
}

// -------------------------------- POST Methods --------------------------------

func (p *Product) Store() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
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

		// validate missing JSON key:values
		validate := validator.New()
		if err := validate.Struct(&request); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error(), "data": nil})
			return
		}
		// Validate Code Value
		if p.s.InvalidCodeValue(request.Code_value) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "there is already a product with that code", "data": nil})
			return
		}

		// Process
		product, err := p.s.Store(request.Name, request.Quantity, request.Code_value, request.Is_published, request.Expiration, request.Price)
		if err != nil {
			ctx.JSON(500, nil)
			return
		}

		// Response
		ctx.JSON(http.StatusCreated, response.Ok("succeed to upload a product", product))
	}
}

// -------------------------------- PUT Methods --------------------------------
func (p *Product) UpdateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		var request Request

		token := ctx.GetHeader("token")
		if token != secretKey {
			ctx.JSON(http.StatusUnauthorized, response.Err(ErrUnauthorized))
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid ID", "data": nil})
			return
		}

		err = ctx.ShouldBindJSON(&request)
		fmt.Println(request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": nil})
			return
		}

		// validate missing JSON key:values
		validate := validator.New()
		if err := validate.Struct(&request); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error(), "data": nil})
			return
		}

		// Validate Code Value
		if p.s.InvalidCodeValue(request.Code_value) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "there is already a product with that code", "data": nil})
			return
		}

		// process
		product, err := p.s.Update(id, request.Name, request.Quantity, request.Code_value, request.Is_published, request.Expiration, request.Price)
		if err != nil {
			ctx.JSON(500, nil)
			return
		}

		// response
		ctx.JSON(http.StatusCreated, response.Ok("succeed to update a product", product))
	}
}

// -------------------------------- PATCH Methods --------------------------------
func (p *Product) UpdateProductName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		var request Request

		token := ctx.GetHeader("token")
		if token != secretKey {
			ctx.JSON(http.StatusUnauthorized, response.Err(ErrUnauthorized))
			return
		}

		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid ID", "data": nil})
			return
		}

		err = ctx.ShouldBindJSON(&request)
		fmt.Println(request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": nil})
			return
		}

		// Validate Name
		if request.Name == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "required product name", "data": nil})
			return
		}
		// process
		product, err := p.s.UpdateName(int(id), request.Name)
		if err != nil {
			ctx.JSON(404, gin.H{"error": err.Error()})
			return
		}

		// response
		ctx.JSON(http.StatusCreated, response.Ok("succeed to update the Name of a product", product))
	}

}

// -------------------------------- DELETE Methods --------------------------------
func (p *Product) DeleteProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Request
		token := ctx.GetHeader("token")
		if token != secretKey {
			ctx.JSON(http.StatusUnauthorized, response.Err(ErrUnauthorized))
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid ID", "data": nil})
			return
		}

		// Process
		err = p.s.Delete(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": nil})
			return
		}

		// Response
		ctx.JSON(http.StatusCreated, response.Ok("succeed to delete the product", id))

	}
}
