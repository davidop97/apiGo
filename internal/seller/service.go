package seller

import (
	"context"
	"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var (
	ErrSellerAlreadyExists = errors.New("seller already exists")
	ErrCannotSaveSeller    = errors.New("cannot add new seller")
	ErrSellerNotExists     = errors.New("seller not exists")
	ErrUpdateSeller        = errors.New("cannot update this seller")
)

type Service interface {
	GetAllSellers(ctx context.Context) ([]domain.Seller, error)
	GetSellerByID(ctx context.Context, id int) (domain.Seller, error)
	Save(ctx context.Context, seller domain.Seller) (int, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, seller domain.Seller, id int) error
	GetLocalityIdFromSeller(ctx context.Context, id int) bool
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

// Get all the sellers available.
func (s service) GetAllSellers(ctx context.Context) ([]domain.Seller, error) {
	allSellers, err := s.r.GetAll(ctx)
	return allSellers, err
}

// Get a seller by id if exists.
func (s service) GetSellerByID(ctx context.Context, id int) (domain.Seller, error) {
	seller, err := s.r.Get(ctx, id)
	return seller, err
}

// Save a new seller or return an error if can not do that action.
func (s *service) Save(ctx context.Context, seller domain.Seller) (int, error) {
	//Check if the seller already exists.
	exist := s.r.Exists(ctx, seller.CID)
	if exist {
		//If the seller already exists, return an error.
		return 0, ErrSellerAlreadyExists
	}
	//Try to save the seller
	newSeller, err := s.r.Save(ctx, seller)
	if err != nil {
		//If an error occurs, return an error to be controlled in the handler.
		return 0, err
	}
	//If everything is ok, return the id of the new seller created.
	return newSeller, nil
}

// Delete a seller using its id. First check if that id
// exists and then delete it or return an error if the seller
// not exist.
func (s *service) Delete(ctx context.Context, id int) error {
	//Check if the seller exists.
	_, err := s.r.Get(ctx, id)
	if err != nil {
		//If the seller not exists, return an error. It can be a not found  or another error (internal error).
		//They will be controlled in the handler.
		return err
	}
	//Try to delete the seller if exists.
	err2 := s.r.Delete(ctx, id)
	if err2 != nil {
		//If an error occurs, return an error to be controlled in the handler (internal error).
		if errors.Is(err2, ErrNotFound) {
			return ErrNotFound //If the seller not exists, return an error.
		}
		//If occurs another error, return it to be controlled in the handler. It can be an internal error.
		return err2
	}
	//If everything is ok, return nil.
	return nil
}

// Update a seller using its id. Returns error otherwise
func (s *service) Update(ctx context.Context, seller domain.Seller, id int) error {
	err := s.r.Update(ctx, seller)
	return err
}

// Check if locality_id exist from seller when trying to create a new seller.
func (s *service) GetLocalityIdFromSeller(ctx context.Context, id int) bool {
	err := s.r.GetLocalityIdFromSeller(ctx, id)
	return err
}
