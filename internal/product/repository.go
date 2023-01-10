package product

import (
	"encoding/json"
	"fmt"
	"os"
)

type Repository interface {
	GetAll() ([]Product, error)
	Store(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (Product, error)
	LastID() (int, error)
	Update(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (Product, error)
	UpdateName(id int, name string) (Product, error)
	Delete(id int) error
}

var Products []Product

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
