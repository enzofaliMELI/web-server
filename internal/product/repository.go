package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/enzofaliMELI/web-server/internal/domain"
	"github.com/enzofaliMELI/web-server/pkg/store"
)

var (
	ErrNotFound = errors.New("item not found")
)

type Repository interface {
	// Read methods
	GetAll() ([]domain.Product, error)
	GetById(id int) (domain.Product, error)
	GetPriceGt(price float64) ([]domain.Product, error)
	LastID() (int, error)
	// Write methods
	Store(name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error)
	Update(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error)
	UpdatePATCH(id int, request domain.PatchRequest) (domain.Product, error)
	Delete(id int) error
	// Validation methods
	InvalidCodeValue(codeVal string, id int) (ok bool)
}

type repository struct {
	db     *[]domain.Product
	lastID int
}

const filename = "../products.json"

func OpenProducts(db *[]domain.Product) (err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("%w. %s", err, "error opening file:")
	}

	err = json.Unmarshal(data, &db)
	if err != nil {
		return fmt.Errorf("%w. %s", err, "error encoding JSON records:")
	}
	return
}

func NewRepository(db *[]domain.Product) Repository {
	return &repository{db: db, lastID: len(*db)}
}

// --------------------------------- Read methods -----------------------------------

func (r *repository) GetAll() ([]domain.Product, error) {
	// Normal Implementation
	// return *r.db, nil

	// Store Implementation #########################################################################################################################
	var product []domain.Product
	storage := store.NewStorage(filename)
	storage.Get(&product)
	return product, nil
}

func (r *repository) LastID() (int, error) {
	return r.lastID, nil
}

func (r *repository) GetById(id int) (domain.Product, error) {
	for _, product := range *r.db {
		if id != 0 && (product.Id == id) {
			return product, nil
		}
	}
	return domain.Product{}, fmt.Errorf("%w. %s", ErrNotFound, "product does not exist")
}

func (r *repository) GetPriceGt(price float64) ([]domain.Product, error) {

	var filtered []domain.Product
	for _, product := range *r.db {
		if price != 0 && (product.Price <= price) {
			continue
		}
		filtered = append(filtered, product)
	}
	return filtered, nil
}

func (r *repository) InvalidCodeValue(codeVal string, id int) (ok bool) {
	for _, product := range *r.db {
		if codeVal == product.Code_value && id != product.Id {
			return true
		}
	}
	return
}

// --------------------------------- Write methods ----------------------------------

func (r *repository) Store(name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error) {
	r.lastID++

	prod := domain.Product{
		Id:           r.lastID,
		Name:         name,
		Quantity:     quantity,
		Code_value:   code_value,
		Is_published: is_published,
		Expiration:   expiration,
		Price:        price,
	}

	*r.db = append(*r.db, prod)

	// Save JSON #########################################################################################################################
	storage := store.NewStorage(filename)
	storage.Set(r.db)

	return prod, nil
}

func (r *repository) Update(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error) {

	newProduct := domain.Product{
		Id:           id,
		Name:         name,
		Quantity:     quantity,
		Code_value:   code_value,
		Is_published: is_published,
		Expiration:   expiration,
		Price:        price,
	}

	update := false

	for i := range *r.db {
		if (*r.db)[i].Id == id {
			(*r.db)[i] = newProduct
			update = true
			break
		}
	}

	if !update {
		return domain.Product{}, fmt.Errorf("%w. %s", ErrNotFound, "product does not exist")
	}

	return newProduct, nil
}

func (r *repository) UpdatePATCH(id int, request domain.PatchRequest) (domain.Product, error) {
	var prod domain.Product

	update := false

	for i := range *r.db {
		if (*r.db)[i].Id == id {
			if request.Name != nil {
				(*r.db)[i].Name = *request.Name
			}
			if request.Quantity != nil {
				(*r.db)[i].Quantity = *request.Quantity
			}
			if request.Code_value != nil {
				(*r.db)[i].Code_value = *request.Code_value
			}
			if request.Is_published != nil {
				(*r.db)[i].Is_published = *request.Is_published
			}
			if request.Expiration != nil {
				(*r.db)[i].Expiration = *request.Expiration
			}
			if request.Price != nil {
				(*r.db)[i].Price = *request.Price
			}

			update = true
			prod = (*r.db)[i]
			break
		}
	}

	if !update {
		return domain.Product{}, fmt.Errorf("%w. %s", ErrNotFound, "product does not exist")
	}

	return prod, nil
}

func (r *repository) Delete(id int) error {
	delete := false
	index := 0

	for i, product := range *r.db {
		if product.Id == id {
			index = i
			delete = true
			break
		}
	}

	if !delete {
		return fmt.Errorf("%w. %s", ErrNotFound, "product does not exist")
	}

	*r.db = append((*r.db)[:index], (*r.db)[index+1:]...)
	return nil
}
