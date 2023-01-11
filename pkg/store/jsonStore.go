package store

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/enzofaliMELI/web-server/internal/domain"
)

type StorageProducts interface {
	Get(pointer *[]domain.Product) error
	Set(pointer *[]domain.Product) error
}

type storageProducts struct {
	file string
}

func NewStorage(file string) StorageProducts {
	return &storageProducts{file: file}
}

func (sp *storageProducts) Get(pointer *[]domain.Product) (err error) {
	// Open
	file, err := os.ReadFile(sp.file)
	if err != nil {
		return fmt.Errorf("%w. %s", err, "error opening file:")
	}

	// Decoder
	err = json.Unmarshal(file, &pointer)
	if err != nil {
		return fmt.Errorf("%w. %s", err, "error encoding JSON records:")
	}
	return
}

func (sp *storageProducts) Set(pointer *[]domain.Product) error {
	// Create
	file, err := os.Create(sp.file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encoder
	data, _ := json.Marshal(pointer)

	// Save
	_, err = io.WriteString(file, string(data))
	if err != nil {
		return err
	}
	return nil
}
