package product

import (
	"time"

	"github.com/enzofaliMELI/web-server/internal/domain"
)

type Service interface {
	GetAll() ([]domain.Product, error)
	Store(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error)
	Update(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error)
	UpdateName(id int, name string) (domain.Product, error)
	Delete(id int) error
}

// -------------------------------- GET Methods --------------------------------

func GetProducts() []domain.Product {
	return Products
}

func GetProductById(id int) []domain.Product {
	products := GetProducts()

	var filtered []domain.Product
	for _, product := range products {
		if id != 0 && (product.Id != id) {
			continue
		}
		filtered = append(filtered, product)
	}
	return filtered
}

func GetProductByPriceGt(price float64) []domain.Product {
	products := GetProducts()

	var filtered []domain.Product
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
