package seller

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"

	"github.com/stretchr/testify/mock"
)

// ServiceMock struct. Mock of the service.
type ServiceMock struct {
	mock.Mock
}

// NewMockService function. Create a new mock of the service.
func NewMockService() *ServiceMock {
	return &ServiceMock{}
}

// GetAllSellers function. Mock of the GetAllSellers function. Get all the sellers available or an error.
func (s *ServiceMock) GetAllSellers(ctx context.Context) (l []domain.Seller, err error) {
	//Get all the sellers.
	args := s.Called(ctx)
	//Return an array of sellers and the error if exists.
	return args.Get(0).([]domain.Seller), args.Error(1)
}

// GetSellerByID function. Mock of the GetSellerByID function. Get a seller by id if exists or an error.
func (s *ServiceMock) GetSellerByID(ctx context.Context, id int) (i domain.Seller, err error) {
	//Get the seller by id.
	args := s.Called(ctx, id)
	//Return the seller and the error if exists.
	return args.Get(0).(domain.Seller), args.Error(1)
}

// Save function. Mock of the Save function. Save a new seller or return an error if can not do that action.
func (s *ServiceMock) Save(ctx context.Context, seller domain.Seller) (id int, err error) {
	//Try to Save the seller.
	args := s.Called(ctx, seller)
	//Return the id of the seller if could be added and the error if exists.
	return args.Int(0), args.Error(1)
}

// Update function. Mock of the Update function. Update a seller or return an error if can not do that action.
func (s *ServiceMock) Update(ctx context.Context, seller domain.Seller, id int) (err error) {
	//Try to Update the seller.
	args := s.Called(ctx, seller)
	//Return the error if exists. Otherwise returns nil.
	return args.Error(0)
}

// Delete function. Mock of the Delete function. Delete a seller or return an error if can not do that action.
func (s *ServiceMock) Delete(ctx context.Context, id int) (err error) {
	//Try to Delete the seller.
	args := s.Called(ctx, id)
	//Return the error if exists. Otherwise returns nil.
	return args.Error(0)
}

// GetLocalityIdFromSeller function. Mock of the GetLocalityIdFromSeller function. Check if the locality id
// from a seller exists. Return true if exists or false if not.
func (s *ServiceMock) GetLocalityIdFromSeller(ctx context.Context, id int) bool {
	//Try to check if the locality id from a seller exists.
	args := s.Called(ctx, id)
	//Return true if exists. Otherwise returns false.
	return args.Bool(0)
}
