package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SERVICE -------------------------------------------------------------------------------------------

type Product struct {
	Id           int
	Name         string
	Quantity     int
	Code_value   string
	Is_published bool
	Expiration   string
	Price        float64
}

/*
var products = []Product{
	{Id: 1, Name: "Oil - Margarine", Quantity: 439, Code_value: "S82254D", Is_published: true, Expiration: "15/12/2021", Price: 71.42},
	{Id: 2, Name: "Pineapple - Canned, Rings", Quantity: 345, Code_value: "M4637", Is_published: true, Expiration: "09/08/2021", Price: 352.79},
}
*/

var products []Product

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

// CONTROLLER -------------------------------------------------------------------------------------------

func Pong(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

func Products(ctx *gin.Context) {
	// process
	products = GetProducts()

	// response
	ctx.JSON(http.StatusOK, gin.H{"message": "succeed to get all websites", "data": products})
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
	ctx.JSON(http.StatusOK, gin.H{"message": "succeed to get all websites", "data": products})
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
	ctx.JSON(http.StatusOK, gin.H{"message": "succeed to get all websites", "data": products})
}

// SERVER -------------------------------------------------------------------------------------------

const filename = "products.json"

func main() {
	OpenProducts(filename)
	server := gin.Default()

	server.GET("/ping", Pong)
	server.GET("/products", Products)
	server.GET("/products/:id", ProductsId)
	server.GET("/products/search", ProductsSearch)

	server.Run(":8080")
}
