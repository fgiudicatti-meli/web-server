package product

import (
	"fmt"
	"github.com/fgiudicatti-meli/web-server/internal/domain"
)

type Service interface {
	GetAll() ([]domain.Product, error)
	Save(name, codeValue, expiration string, quantity int, price float64, isPublished bool) (domain.Product, error)
	GetById(id int) (domain.Product, error)
	Update(id int, name, codeValue, expiration string, quantity int, price float64, isPublished bool) (domain.Product, error)
	Delete(id int) error
	UpdateName(id int, name string) (domain.Product, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) GetAll() ([]domain.Product, error) {
	allProducts, err := s.repository.GetAll()
	if err != nil {
		return []domain.Product{}, fmt.Errorf("la lista de productos esta vacia")
	}

	return allProducts, nil
}

func (s *service) Save(name, codeValue, expiration string, quantity int, price float64, isPublished bool) (domain.Product, error) {
	lastId, err := s.repository.GetLastId()
	if err != nil {
		return domain.Product{}, fmt.Errorf("internal server error")
	}

	lastId++

	prd, err := s.repository.Save(lastId, name, codeValue, expiration, quantity, price, isPublished)
	if err != nil {
		return domain.Product{}, fmt.Errorf("algo paso al intentar guardar controle el cuerpo y el id ingresado")
	}

	return prd, nil
}

func (s *service) GetById(id int) (domain.Product, error) {
	productById, err := s.repository.GetById(id)
	if err != nil {
		return domain.Product{}, fmt.Errorf("no se encontro el producto con %d", id)
	}

	return productById, nil
}

func (s *service) Update(id int, name, codeValue, expiration string, quantity int, price float64, isPublished bool) (domain.Product, error) {
	return s.repository.Update(id, name, codeValue, expiration, quantity, price, isPublished)
}

func (s *service) Delete(id int) error {
	return s.repository.Delete(id)
}

func (s *service) UpdateName(id int, name string) (domain.Product, error) {
	return s.repository.UpdateName(id, name)
}
