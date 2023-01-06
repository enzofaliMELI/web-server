package services

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/enzofaliMELI/web-server/internal/product/models"
)

var Products []models.Product

var LastID int

func OpenProducts(filename string) (err error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	err = json.Unmarshal(data, &Products)
	if err != nil {
		fmt.Println("Error encoding json records:", err)
		return
	}

	LastID = len(Products)
	return
}

func GetProducts() []models.Product {
	return Products
}

func GetProductById(id int) []models.Product {
	products := GetProducts()

	var filtered []models.Product
	for _, product := range products {
		if id != 0 && (product.Id != id) {
			continue
		}
		filtered = append(filtered, product)
	}
	return filtered
}

func GetProductByPriceGt(price float64) []models.Product {
	products := GetProducts()

	var filtered []models.Product
	for _, product := range products {
		if price != 0 && (product.Price <= price) {
			continue
		}
		filtered = append(filtered, product)
	}
	return filtered
}

/*
func InvalidJSONBody(newProduct models.Product) (ok bool) {
	if newProduct.Name == "" || newProduct.Quantity == 0 || newProduct.Code_value == "" || newProduct.Expiration == "" || newProduct.Price == 0 {
		return true
	}
	return
}
*/

func InvalidCodeValue(newProduct models.Product, products []models.Product) (ok bool) {
	for _, p := range products {
		if newProduct.Code_value == p.Code_value {
			return true
		}
	}
	return
}

func InvalidExpiration(newProduct models.Product) (ok bool) {
	// Parse the date using the layout "02/01/2006"
	_, err := time.Parse("02/01/2006", newProduct.Expiration)
	if err != nil {
		return true
	}
	return
}

func CreateProducts() []models.Product {
	return Products
}
