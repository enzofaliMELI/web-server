package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/enzofaliMELI/web-server/internal/domain"
	"github.com/enzofaliMELI/web-server/internal/product"
	"github.com/enzofaliMELI/web-server/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var (
	ErrUnauthorized      = errors.New("error: invalid token")
	ErrInvalidId         = errors.New("error: invalid Id")
	ErrInvalidCodeValue  = errors.New("invalid expiration date, (format: DD/MM/YYYY)")
	ErrInvalidExpiration = errors.New("there is already a product with that code")
)

type Product struct {
	s product.Service
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
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}

		// Response
		ctx.JSON(http.StatusOK, response.Ok("succeed to get all products", products))
	}
}

func (p *Product) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Err(ErrInvalidId))
			return
		}
		// Process
		products, err := p.s.GetById(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}

		// Response
		ctx.JSON(http.StatusOK, response.Ok("succeed to get product", products))
	}
}

func (p *Product) GetPriceGt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		price, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Err(ErrInvalidId))
			return
		}
		// Process
		products, err := p.s.GetPriceGt(price)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}

		// Response
		ctx.JSON(http.StatusOK, response.Ok("succeed to get all products", products))
	}
}

// -------------------------------- POST Methods --------------------------------

func (p *Product) Store() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		var request domain.Request

		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, response.Err(ErrUnauthorized))
			return
		}

		err := ctx.ShouldBindJSON(&request)
		fmt.Println(request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Err(err))
			return
		}

		// validate missing JSON key:values
		validate := validator.New()
		if err := validate.Struct(&request); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, response.Err(err))
			return
		}

		// Validate Code Value
		if p.s.InvalidCodeValue(request.Code_value, -1) {
			ctx.JSON(http.StatusBadRequest, response.Err(ErrInvalidCodeValue))
			return
		}

		// Validate Expiration date format
		if p.s.InvalidExpiration(request.Expiration) {
			ctx.JSON(http.StatusBadRequest, response.Err(ErrInvalidExpiration))
			return
		}

		// Process
		product, err := p.s.Store(request.Name, request.Quantity, request.Code_value, request.Is_published, request.Expiration, request.Price)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
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
		var request domain.Request

		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, response.Err(ErrUnauthorized))
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Err(ErrInvalidId))
			return
		}

		err = ctx.ShouldBindJSON(&request)
		fmt.Println(request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Err(err))
			return
		}

		// validate missing JSON key:values
		validate := validator.New()
		if err := validate.Struct(&request); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, response.Err(err))
			return
		}

		// Validate Code Value
		if p.s.InvalidCodeValue(request.Code_value, id) {
			ctx.JSON(http.StatusBadRequest, response.Err(ErrInvalidCodeValue))
			return
		}

		// Validate Expiration date format
		if p.s.InvalidExpiration(request.Expiration) {
			ctx.JSON(http.StatusBadRequest, response.Err(ErrInvalidExpiration))
			return
		}

		// Process
		product, err := p.s.Update(id, request.Name, request.Quantity, request.Code_value, request.Is_published, request.Expiration, request.Price)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}

		// Response
		ctx.JSON(http.StatusCreated, response.Ok("succeed to update a product", product))
	}
}

// -------------------------------- PATCH Methods --------------------------------
func (p *Product) UpdatePATCH() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Request
		var request domain.PatchRequest

		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, response.Err(ErrUnauthorized))
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Err(ErrInvalidId))
			return
		}

		err = ctx.ShouldBindJSON(&request)
		fmt.Println(request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Err(err))
			return
		}

		// Validate Code Value
		if request.Code_value != nil {
			if p.s.InvalidCodeValue(*request.Code_value, id) {
				ctx.JSON(http.StatusBadRequest, response.Err(ErrInvalidCodeValue))
				return
			}
		}

		// Process
		product, err := p.s.UpdatePATCH(id, request)
		if err != nil {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		// Response
		ctx.JSON(http.StatusCreated, response.Ok("succeed to update the Name of a product", product))
	}

}

// -------------------------------- DELETE Methods --------------------------------
func (p *Product) DeleteProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Request
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, response.Err(ErrUnauthorized))
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.Err(ErrInvalidId))
			return
		}

		// Process
		err = p.s.Delete(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}

		// Response
		ctx.JSON(http.StatusCreated, response.Ok("succeed to delete the product", id))

	}
}
