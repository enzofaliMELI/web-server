package repository

import (
	"encoding/json"
	"fmt"
	"os"

	service "github.com/enzofaliMELI/web-server/internal/product/service"
)

func OpenProducts(filename string) (err error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	err = json.Unmarshal(data, &service.Products)
	if err != nil {
		fmt.Println("Error encoding json records:", err)
		return
	}

	service.LastID = len(service.Products)
	return
}
