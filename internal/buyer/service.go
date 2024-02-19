package buyer

import (
	"context"
	"errors"
	"fmt"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var (
	ErrNotFound      = errors.New("buyer not found")
	ErrAlreadyExists = errors.New("buyer already exists")
)

// Service is an interface that defines methods for a service
type Service interface {
	// GetAll obtains all buyers
	GetAll(ctx context.Context) ([]domain.Buyer, error)
	// Get obtains a buyer by id
	Get(ctx context.Context, id int) (domain.Buyer, error)
	// Delete a buyer by id
	Delete(ctx context.Context, id int) error
	// Save adds a new buyer
	Save(ctx context.Context, b domain.Buyer) (int, error)
	// Update a buyer by id
	Update(ctx context.Context, id int, b domain.Buyer, bs *domain.Buyer) error
}

// service is the concrete implementation of the service interface
type service struct {
	r Repository
}

// NewService creates a new service
func NewService(r Repository) Service {
	return &service{r}
}

// GetAll returns all buyers
func (s *service) GetAll(ctx context.Context) ([]domain.Buyer, error) {
	l, err := s.r.GetAll(ctx)
	return l, err
}

// Get returns a Buyer by id
func (s *service) Get(ctx context.Context, id int) (domain.Buyer, error) {
	// call to Get function to obtain a buyer if it exist
	b, err := s.r.Get(ctx, id)
	// check if buyer not found
	if err != nil {
		return domain.Buyer{}, ErrNotFound
	}
	// return the buyer founded
	return b, nil
}

// Get a buyer by id and delete it if there
func (s *service) Delete(ctx context.Context, id int) error {
	_, err := s.r.Get(ctx, id)
	// check if buyer not found
	if err != nil {
		return ErrNotFound
	}
	// if it was founded, delete it.
	err = s.r.Delete(ctx, id)
	if err != nil {
		return err
	}
	// return nil if there´s no error
	return nil
}

// Save a new buyer to the database
func (s *service) Save(ctx context.Context, b domain.Buyer) (int, error) {
	// check if card number already exists in the database
	fmt.Printf(" card number id: %v", b.CardNumberID)
	exists := s.r.Exists(ctx, b.CardNumberID)
	// check if buyer already exists
	if exists {
		return 0, ErrAlreadyExists
	}
	// save the buyer
	id, err := s.r.Save(ctx, b)
	if err != nil {
		return 0, err
	}
	// return id of the buyer saved into the database
	return id, nil
}

// Update buyer information if the buyer exists in the database
func (s *service) Update(ctx context.Context, id int, u domain.Buyer, b *domain.Buyer) error {

	// check if buyer already exists in the database
	exists := s.r.Exists(ctx, u.CardNumberID)
	if exists {
		return ErrAlreadyExists
	}
	// check if fields are correct
	if u.FirstName != "" {
		b.FirstName = u.FirstName
	}
	if u.LastName != "" {
		b.LastName = u.LastName
	}
	// update the buyer
	err := s.r.Update(ctx, *b)
	if err != nil {
		return err
	}
	// return nil if there no error
	return nil
}
