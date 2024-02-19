package product

import (
	"context"
	"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var (
	ErrNotFound          = errors.New("product not found")
	ErrProductCodeExists = errors.New("product_code already exists")
	ErrorSavingProduct   = errors.New("error saving product")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	Get(ctx context.Context, id int) (domain.Product, error)
	Exists(ctx context.Context, productCode string) bool
	Save(ctx context.Context, p domain.Product) (int, error)
	Update(ctx context.Context, p domain.Product) error
	Delete(ctx context.Context, id int) error
	CreateProductRecord(ctx context.Context, p domain.ProductRecordCreate) (int, error)
	GetProductRecord(ctx context.Context, idProduct int) ([]domain.ProductRecordGet, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetAll retrieves all products from the database.
// It returns a slice of domain.Product and an error if there is any.
func (s *service) GetAll(ctx context.Context) ([]domain.Product, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, ErrNotFound
	}
	return products, nil

}

// Get retrieves a product by its ID from the database.
// It returns the domain.Product and an error if there is any.
func (s *service) Get(ctx context.Context, id int) (domain.Product, error) {
	product, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Product{}, ErrNotFound
	}
	return product, nil
}

// Exists checks if a product with the given productCode exists in the database.
// It returns a boolean indicating the existence of the product.
func (s *service) Exists(ctx context.Context, productCode string) bool {
	return s.repo.Exists(ctx, productCode)
}

// Save saves a new product in the database.
// It returns the ID of the saved product and an error if there is any.
func (s *service) Save(ctx context.Context, p domain.Product) (int, error) {
	// check if product_code already exists
	if s.Exists(ctx, p.ProductCode) {
		return 0, ErrProductCodeExists
	}

	product, err := s.repo.Save(ctx, p)
	if err != nil {
		return 0, ErrorSavingProduct
	}
	return product, nil
}

// Update updates an existing product in the database.
// It returns an error if there is any.
func (s *service) Update(ctx context.Context, p domain.Product) error {
	existProduct, err := s.repo.Get(ctx, p.ID)
	if err != nil {
		return err
	}

	if existProduct.ProductCode != p.ProductCode {
		if s.repo.Exists(ctx, p.ProductCode) {
			return ErrProductCodeExists
		}
	}

	err = s.repo.Update(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes an existing product from the database.
// It returns an error if there is any.
func (s *service) Delete(ctx context.Context, id int) error {
	_, err := s.repo.Get(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// CreateProductRecord creates a new product record in the database.
// It returns the ID of the created product record and an error if there is any.
func (s *service) CreateProductRecord(ctx context.Context, p domain.ProductRecordCreate) (int, error) {
	//check if product exists in table products
	_, err := s.repo.Get(ctx, p.ProductID)
	if err != nil {
		return 0, ErrNotFound
	}

	product, err := s.repo.CreateProductRecord(ctx, p)
	if err != nil {
		return 0, err
	}
	return product, nil
}

// GetProductRecord retrieves product records by product ID from the database.
// If idProduct is 0, it retrieves all product records.
// It returns a slice of domain.ProductRecordGet and an error if there is any.
func (s *service) GetProductRecord(ctx context.Context, idProduct int) ([]domain.ProductRecordGet, error) {

	// Check, if the id is different from 0, check if the product exists in the table products
	if idProduct != 0 {
		_, err := s.repo.Get(ctx, idProduct) //Check if exist
		if err != nil {
			return nil, ErrNotFound
		}
	}

	product, err := s.repo.GetProductRecord(ctx, idProduct)
	if err != nil {
		return nil, err
	}
	return product, nil
}
