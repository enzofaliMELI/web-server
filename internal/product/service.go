package product

import (
	"time"
)

type Service interface {
	GetAll() ([]Product, error)
	Store(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (Product, error)
	Update(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (Product, error)
	UpdateName(id int, name string) (Product, error)
	Delete(id int) error
}

// -------------------------------- GET Methods --------------------------------

func GetProducts() []Product {
	return Products
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

// -------------------------------- Validation Methods --------------------------------

func InvalidCodeValue(codeVal string) (ok bool) {
	for _, p := range Products {
		if codeVal == p.Code_value {
			return true
		}
	}
	return
}

func InvalidExpiration(expiration string) (ok bool) {
	// Parse the date using the layout "02/01/2006"
	_, err := time.Parse("02/01/2006", expiration)
	if err != nil {
		return true
	}
	return
}
