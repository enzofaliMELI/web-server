package product

import (
	"time"

	"github.com/enzofaliMELI/web-server/internal/domain"
)

type Service interface {
	// Read methods
	GetAll() ([]domain.Product, error)
	GetById(id int) (domain.Product, error)
	GetPriceGt(price float64) ([]domain.Product, error)
	// Write methods
	Store(name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error)
	Update(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error)
	UpdatePATCH(id int, request domain.PatchRequest) (domain.Product, error)
	Delete(id int) error
	// Validation
	InvalidCodeValue(codeVal string, id int) (ok bool)
	InvalidExpiration(expiration string) (ok bool)
}

func NewService(r Repository) Service {
	return &service{r: r}
}

type service struct {
	// repo
	r Repository

	// external api's
	// ...
}

// --------------------------------- Read methods -----------------------------------

func (s *service) GetAll() ([]domain.Product, error) {
	return s.r.GetAll()
}

func (s *service) GetById(id int) (domain.Product, error) {
	return s.r.GetById(id)
}

func (s *service) GetPriceGt(price float64) ([]domain.Product, error) {
	return s.r.GetPriceGt(price)
}

// --------------------------------- Write methods ----------------------------------

func (s *service) Store(name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error) {

	prod, err := s.r.Store(name, quantity, code_value, is_published, expiration, price)
	if err != nil {
		return domain.Product{}, err
	}
	return prod, nil

}

func (s *service) Update(id int, name string, quantity int, code_value string, is_published bool, expiration string, price float64) (domain.Product, error) {
	return s.r.Update(id, name, quantity, code_value, is_published, expiration, price)
}

func (s *service) UpdatePATCH(id int, request domain.PatchRequest) (domain.Product, error) {
	return s.r.UpdatePATCH(id, request)
}

func (s *service) Delete(id int) error {
	return s.r.Delete(id)
}

// -------------------------------- Validation Methods --------------------------------

func (s *service) InvalidCodeValue(codeVal string, id int) (ok bool) {
	return s.r.InvalidCodeValue(codeVal, id)
}

func (r *service) InvalidExpiration(expiration string) (ok bool) {
	// Parse the date using the layout "02/01/2006"
	_, err := time.Parse("02/01/2006", expiration)
	if err != nil {
		return true
	}
	return
}
