package batch

import (
	"context"
	"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var (
	ErrDuplicateBatchNumber = errors.New("duplicate batch number")
)

type Service interface {
	GetAll(ctx context.Context) (l []domain.ProductBatch, err error)
	Save(ctx context.Context, batch domain.ProductBatch) (id int, err error)
}

// service is a struct that represents a ProductBatch service
type service struct {
	r Repository
}

// NewService returns a new instance of ProductBatch service
func NewService(r Repository) Service {
	return &service{r}
}

// GetAll returns all Product Batches
func (s *service) GetAll(ctx context.Context) (l []domain.ProductBatch, err error) {
	l, err = s.r.GetAll(ctx)
	return
}

// Save stores a new Product Batch
func (s *service) Save(ctx context.Context, b domain.ProductBatch) (id int, err error) {
	// Check if batch number is unique
	exists := s.r.Exists(ctx, b.BatchNumber)
	if exists {
		err = ErrDuplicateBatchNumber
		return
	}

	// Save new product batch
	id, err = s.r.Save(ctx, b)
	return
}
