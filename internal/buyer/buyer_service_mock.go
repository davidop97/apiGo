package buyer

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/stretchr/testify/mock"
)

// BuyerServiceMock is a struct that represent service  mock
type BuyerServiceMock struct {
	mock.Mock
}

// NewBuyerService creates a new mock of the buyer service
func NewBuyerService() *BuyerServiceMock {
	return &BuyerServiceMock{}
}

// GetAll returns all buyers
func (b *BuyerServiceMock) GetAll(ctx context.Context) ([]domain.Buyer, error) {
	args := b.Called(ctx)
	return args.Get(0).([]domain.Buyer), args.Error(1)
}

// Get returns a buyer by id and error if any
func (b *BuyerServiceMock) Get(ctx context.Context, id int) (domain.Buyer, error) {
	args := b.Called(ctx, id)
	return args.Get(0).(domain.Buyer), args.Error(1)
}

// Delete function get a buyer by id and delete if there is
func (b *BuyerServiceMock) Delete(ctx context.Context, id int) error {
	args := b.Called(ctx, id)
	return args.Error(0)
}

// Save save the buyer to the database
func (b *BuyerServiceMock) Save(ctx context.Context, bu domain.Buyer) (int, error) {
	args := b.Called(ctx, bu)
	return args.Int(0), args.Error(1)
}

// Update update buyer information if the buyer exists in the database
func (b *BuyerServiceMock) Update(ctx context.Context, id int, buyerToUpdate domain.Buyer, buyerInDatabase *domain.Buyer) error {
	args := b.Called(ctx, id, buyerToUpdate, buyerInDatabase)
	*buyerInDatabase = buyerToUpdate // Asegúrese de que buyerInDatabase se actualice con buyerToUpdate.
	return args.Error(0)
}
