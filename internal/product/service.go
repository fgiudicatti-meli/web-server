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
		return nil, err
	}

	return allProducts, nil
}

func (s *service) Save(name, codeValue, expiration string, quantity int, price float64, isPublished bool) (domain.Product, error) {
	lastId, err := s.repository.GetLastId()
	if err != nil {
		return domain.Product{}, fmt.Errorf("error in generate last id: %w", err)
	}

	lastId++

	newProduct, err := s.repository.Save(lastId, name, codeValue, expiration, quantity, price, isPublished)
	if err != nil {
		return domain.Product{}, fmt.Errorf("error when try creating product: %w", err)
	}

	return newProduct, nil
}

func (s *service) GetById(id int) (domain.Product, error) {
	productById, err := s.repository.GetById(id)
	if err != nil {
		return domain.Product{}, err
	}

	return productById, nil
}

func (s *service) Update(id int, name, codeValue, expiration string, quantity int, price float64, isPublished bool) (domain.Product, error) {
	product, err := s.repository.Update(id, name, codeValue, expiration, quantity, price, isPublished)
	if err != nil {
		return domain.Product{}, fmt.Errorf("error when try update user: %w", err)
	}
	return product, nil
}

func (s *service) Delete(id int) error {
	err := s.repository.Delete(id)
	if err != nil {
		return fmt.Errorf("error when try deleted product: %w", err)
	}

	return nil
}

func (s *service) UpdateName(id int, name string) (domain.Product, error) {
	product, err := s.repository.UpdateName(id, name)
	if err != nil {
		return domain.Product{}, fmt.Errorf("error when update user: %w", err)
	}
	return product, nil
}
