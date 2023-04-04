package product

import (
	"fmt"
	"github.com/fgiudicatti-meli/web-server/internal/domain"
	"github.com/fgiudicatti-meli/web-server/pkg/store"
)

const (
	ErrProductNotFound     = "product with id: %d not found"
	ErrFailWhenReadingFile = "fail when try read file"
	ErrFailWhenWritingFile = "can't write in file"
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

type repository struct {
	db store.Store
}

func NewRepository(db store.Store) Repository {
	return &repository{
		db: db,
	}
}

// var products []domain.Product
// var lastId = products[len(products)-1].Id
/*
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
*/

func (r *repository) GetAll() ([]domain.Product, error) {
	var users []domain.Product
	if err := r.db.Read(&users); err != nil {
		return nil, fmt.Errorf(ErrFailWhenReadingFile)
	}
	return users, nil
}

func (r *repository) Save(id int, name, codeValue, expiration string, quantity int, price float64, isPublished bool) (domain.Product, error) {
	var products []domain.Product

	if err := r.db.Read(&products); err != nil {
		return domain.Product{}, fmt.Errorf(ErrFailWhenReadingFile)
	}

	newProduct := domain.Product{
		Id:          id,
		Name:        name,
		Quantity:    quantity,
		CodeValue:   codeValue,
		IsPublished: isPublished,
		Expiration:  expiration,
		Price:       price,
	}
	products = append(products, newProduct)

	if err := r.db.Write(&products); err != nil {
		return domain.Product{}, fmt.Errorf(ErrFailWhenWritingFile)
	}

	return newProduct, nil
}

func (r *repository) GetLastId() (int, error) {
	var products []domain.Product
	if err := r.db.Read(&products); err != nil {
		return 0, fmt.Errorf(ErrFailWhenReadingFile)
	}

	if len(products) == 0 {
		return 0, nil
	}
	return products[len(products)-1].Id, nil
}

func (r *repository) GetById(id int) (domain.Product, error) {
	var products []domain.Product
	if err := r.db.Read(&products); err != nil {
		return domain.Product{}, fmt.Errorf(ErrFailWhenReadingFile)
	}

	var targetProduct domain.Product
	founded := false
	for i := range products {
		if products[i].Id == id {
			targetProduct = products[i]
			founded = true
			break
		}
	}

	if !founded {
		return domain.Product{}, fmt.Errorf(ErrProductNotFound, id)
	}

	return targetProduct, nil
}

func (r *repository) Update(id int, name, codeValue, expiration string, quantity int, price float64, isPublished bool) (domain.Product, error) {
	var products []domain.Product
	if err := r.db.Read(&products); err != nil {
		return domain.Product{}, fmt.Errorf(ErrFailWhenReadingFile)
	}

	updateUser := domain.Product{
		Name:        name,
		Quantity:    quantity,
		CodeValue:   codeValue,
		IsPublished: isPublished,
		Expiration:  expiration,
		Price:       price,
	}
	hasChange := false

	for i := range products {
		if products[i].Id == id {
			updateUser.Id = id
			products[i] = updateUser
			hasChange = true
			break
		}
	}

	if !hasChange {
		return domain.Product{}, fmt.Errorf(ErrProductNotFound, id)
	}

	if err := r.db.Write(&products); err != nil {
		return domain.Product{}, fmt.Errorf(ErrFailWhenWritingFile)
	}

	return updateUser, nil
}

func (r *repository) UpdateName(id int, name string) (domain.Product, error) {
	var products []domain.Product
	if err := r.db.Read(&products); err != nil {
		return domain.Product{}, fmt.Errorf(ErrFailWhenReadingFile)
	}

	hasChange := false
	var updateProduct domain.Product
	for i := range products {
		if products[i].Id == id {
			products[i].Name = name
			hasChange = true
			updateProduct = products[i]
			break
		}
	}

	if !hasChange {
		return domain.Product{}, fmt.Errorf(ErrProductNotFound, id)
	}

	if err := r.db.Write(&products); err != nil {
		return domain.Product{}, fmt.Errorf(ErrFailWhenWritingFile)
	}

	return updateProduct, nil
}

func (r *repository) Delete(id int) (err error) {
	var products []domain.Product
	if err := r.db.Read(&products); err != nil {
		return fmt.Errorf(ErrFailWhenReadingFile)
	}

	var indexProductToDelete int
	hasFounded := false
	for i := range products {
		if products[i].Id == id {
			indexProductToDelete = i
			hasFounded = true
			break
		}
	}

	if !hasFounded {
		return fmt.Errorf(ErrProductNotFound, id)
	}

	products = append(products[:indexProductToDelete], products[indexProductToDelete+1:]...)
	if err := r.db.Write(&products); err != nil {
		return fmt.Errorf(ErrFailWhenWritingFile)
	}

	return nil
}
