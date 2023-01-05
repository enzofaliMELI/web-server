package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SERVICE -------------------------------------------------------------------------------------------

type Product struct {
	Id           int     `json:"id,omitempty"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

/*
var products = []Product{
	{Id: 1, Name: "Oil - Margarine", Quantity: 439, Code_value: "S82254D", Is_published: true, Expiration: "15/12/2021", Price: 71.42},
	{Id: 2, Name: "Pineapple - Canned, Rings", Quantity: 345, Code_value: "M4637", Is_published: true, Expiration: "09/08/2021", Price: 352.79},
}
*/

var products []Product

var lastID int

func OpenProducts(filename string) (err error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	err = json.Unmarshal(data, &products)
	if err != nil {
		fmt.Println("Error encoding json records:", err)
		return
	}

	lastID = len(products)
	return
}

func GetProducts() []Product {
	return products
}

func GetProductById(id int) []Product {
	products := GetProducts()

	var filtered []Product
	for _, product := range products {
		if id != 0 && (product.Id != id) {
			continue
		}
		filtered = append(filtered, product)
	}
	return filtered
}

func GetProductByPriceGt(price float64) []Product {
	products := GetProducts()

	var filtered []Product
	for _, product := range products {
		if price != 0 && (product.Price <= price) {
			continue
		}
		filtered = append(filtered, product)
	}
	return filtered
}

func InvalidJSONBody(newProduct Product) (ok bool) {
	if newProduct.Name == "" || newProduct.Quantity == 0 || newProduct.Code_value == "" || newProduct.Expiration == "" || newProduct.Price == 0 {
		return true
	}
	return
}

func InvalidCodeValue(newProduct Product, products []Product) (ok bool) {
	for _, p := range products {
		if newProduct.Code_value == p.Code_value {
			return true
		}
	}
	return
}

func InvalidExpiration(newProduct Product) (ok bool) {
	_, err := time.Parse("02/01/2006", newProduct.Expiration)
	if err != nil {
		return true
	}
	return
}

// CONTROLLER -------------------------------------------------------------------------------------------

// GET

func Pong(ctx *gin.Context) {
	// response
	ctx.String(http.StatusOK, "pong")
}

func GETProducts(ctx *gin.Context) {
	// process
	products = GetProducts()

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
	products := GetProductById(id)

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
	products := GetProductByPriceGt(price)

	// response
	ctx.JSON(http.StatusOK, gin.H{"message": "succeed to get all products", "data": products})
}

// POST

func POSTProducts(ctx *gin.Context) {
	// request
	var request Product

	/*
		token := ctx.GetHeader("token")
		if token != "1234" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token", "data": nil})
			return
		}
	*/

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "data": nil})
		return
	}

	//Todo: si Is_published == "" se debe bindear con JSON a false

	if InvalidJSONBody(request) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON, (only id and is_published can be empty)", "data": nil})
		return
	}

	if InvalidCodeValue(request, products) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "there is already a product with that code", "data": nil})
		return
	}

	if InvalidExpiration(request) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid expiration date, (format: DD/MM/YYYY)", "data": nil})
		return
	}

	// process

	lastID++
	request.Id = lastID
	products = append(products, request)

	// response
	ctx.JSON(http.StatusOK, gin.H{"message": "succeed to upload a product", "data": request})
}

// SERVER -------------------------------------------------------------------------------------------

const filename = "products.json"

func main() {
	OpenProducts(filename)
	server := gin.Default()

	server.GET("/ping", Pong)

	products := server.Group("/products")
	products.POST("/", POSTProducts)
	products.GET("/", GETProducts)
	products.GET("/:id", ProductsId)
	products.GET("/search", ProductsSearch)

	server.Run(":8080")
}
