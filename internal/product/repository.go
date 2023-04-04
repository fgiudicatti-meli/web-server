package product

import (
	"encoding/json"
	"fmt"
	"github.com/fgiudicatti-meli/web-server/internal/domain"
	"log"
	"os"
)

const (
	ErrProductNotFound = "user %d not found"
)

type Repository interface {
	GetAll() ([]domain.Product, error)
	Save(id int, name, codeValue, expiration string, quantity int, price float64, isPublished bool) (domain.Product, error)
	GetLastId() (int, error)
	GetById(id int) (domain.Product, error)
	Update(id int, name, codeValue, expiration string, quantity int, price float64, isPublished bool) (domain.Product, error)
	Delete(id int) error
	UpdateName(id int, name string) (domain.Product, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

// var products []domain.Product
// var lastId int
var products = InitValues()
var lastId = products[len(products)-1].Id

func InitValues() []domain.Product {

	productsFromJson, err := os.ReadFile("./products.json")
	if err != nil {
		log.Fatal(err)
	}

	var products []domain.Product
	if err := json.Unmarshal(productsFromJson, &products); err != nil {
		log.Fatal(err)
	}
	return products
}
func (r *repository) GetAll() ([]domain.Product, error) {
	return products, nil
}

func (r *repository) Save(id int, name, codeValue, expiration string, quantity int, price float64, isPublished bool) (domain.Product, error) {
	prd := domain.Product{
		Id:          id,
		Name:        name,
		Quantity:    quantity,
		CodeValue:   codeValue,
		IsPublished: isPublished,
		Expiration:  expiration,
		Price:       price,
	}
	products = append(products, prd)
	lastId = prd.Id
	return prd, nil
}

func (r *repository) GetLastId() (int, error) {
	return lastId, nil
}

func (r *repository) GetById(id int) (domain.Product, error) {
	var target domain.Product
	for _, p := range products {
		if p.Id == id {
			target = p
			return target, nil
		}
	}

	return domain.Product{}, fmt.Errorf(ErrProductNotFound, id)
}

func (r *repository) Update(id int, name, codeValue, expiration string, quantity int, price float64, isPublished bool) (domain.Product, error) {
	p := domain.Product{
		Name:        name,
		Quantity:    quantity,
		CodeValue:   codeValue,
		IsPublished: isPublished,
		Expiration:  expiration,
		Price:       price,
	}

	updated := false
	for i := range products {
		if products[i].Id == id {
			p.Id = id
			products[i] = p
			updated = true
		}
	}

	if !updated {
		return domain.Product{}, fmt.Errorf(ErrProductNotFound, id)
	}
	return p, nil

}

func (r *repository) UpdateName(id int, name string) (domain.Product, error) {
	var productUpdate domain.Product
	updated := false
	for i := range products {
		if products[i].Id == id {
			products[i].Name = name
			updated = true
			productUpdate = products[i]
		}
	}

	if !updated {
		return domain.Product{}, fmt.Errorf(ErrProductNotFound, id)
	}
	return productUpdate, nil
}

func (r *repository) Delete(id int) (err error) {
	deleted := false
	var index int
	for i := range products {
		if products[i].Id == id {
			index = i
			deleted = true
		}
	}

	if !deleted {
		return fmt.Errorf("no existe el usuario con %d", id)
	}

	products = append(products[:index], products[index+1:]...)
	return nil
}
