package product

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/enzofaliMELI/web-server/internal/domain"
)

type Repository interface {
	GetAll() ([]domain.Product, error)
	Store(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error)
	LastID() (int, error)
	Update(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error)
	UpdateName(id int, name string) (domain.Product, error)
	Delete(id int) error
}

var Products []domain.Product

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

func SaveProduct(name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error) {

	prod := domain.Product{
		Id:           0,
		Name:         name,
		Quantity:     quantity,
		Code_value:   code_value,
		Is_published: is_published,
		Expiration:   expiration,
		Price:        price,
	}

	LastID++
	prod.Id = LastID
	Products = append(Products, prod)
	return prod, nil
}
