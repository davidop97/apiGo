package warehouse

import (
	"context"
	"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var (
	ErrNotFound           = errors.New("warehouse not found")
	ErrIncorrectData      = errors.New("incorrect data")
	ErrDuplicateWarehouse = errors.New("warehouse already exists")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Warehouse, error)
	Get(ctx context.Context, id int) (domain.Warehouse, error)
	Save(ctx context.Context, w domain.Warehouse) (int, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, w domain.Warehouse) error
}

type service struct {
	rp Repository
}

func NewService(r Repository) Service {
	return &service{rp: r}
}

// GetAll returns all the warehouses available
func (s *service) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	warehouses, err := s.rp.GetAll(ctx)
	if err != nil {
		return warehouses, ErrNotFound
	}

	return warehouses, nil
}

// Get returns a warehouse by ID, returns error if the warehouse doesn't exists
func (s *service) Get(ctx context.Context, id int) (domain.Warehouse, error) {
	warehouse, err := s.rp.Get(ctx, id)
	if err != nil {
		return warehouse, ErrNotFound
	}

	return warehouse, nil
}

// Save saves a warehouse in the database, returns error if the warehouse already exists or if the data is incorrect
func (s *service) Save(ctx context.Context, w domain.Warehouse) (int, error) {
	//Check for empty strings or negative capacity
	for _, v := range []string{w.Address, w.Telephone, w.WarehouseCode} {
		if v == "" {
			return 0, ErrIncorrectData
		}
	}
	if w.MinimumCapacity < 0 {
		return 0, ErrIncorrectData
	}

	id, err := s.rp.Save(ctx, w)
	if err != nil {
		return id, err
	}

	return id, nil
}

// Delete deletes a warehouse by ID, returns error if the warehouse doesn't exists
func (s *service) Delete(ctx context.Context, id int) error {
	err := s.rp.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// Update updates a warehouse in the database, returns error if the warehouse doesn't exists or if the data is incorrect
func (s *service) Update(ctx context.Context, w domain.Warehouse) error {
	//Check for empty strings or negative capacity
	for _, v := range []string{w.Address, w.Telephone, w.WarehouseCode} {
		if v == "" {
			return ErrIncorrectData
		}
	}
	if w.MinimumCapacity < 0 {
		return ErrIncorrectData
	}

	err := s.rp.Update(ctx, w)
	if err != nil {
		return err
	}

	return nil
}
