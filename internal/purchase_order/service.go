package purchase_order

import (
	"context"
	"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var (
	ErrPurchaseOrderAlreadyExists = errors.New("Purchase order already exists")
	ErrBuyerIDNotExists           = errors.New("Buyer ID not exists")
	ErrProductsRecordIDNotExits   = errors.New("Product record ID not exists")
)

// Service is an interface that defines methods for a service
// that handles purchase order.
type Service interface {
	Save(ctx context.Context, purchaseOrder domain.PurchaseOrder) (int, error)
	PurchaseOrdersByBuyer(ctx context.Context, buyerID int) ([]domain.PurchaseOrdersByBuyer, error)
}

// service struct is the concrete implementation of the Service interface
type service struct {
	repo Repository
}

// NewServices create a new service
func NewService(r Repository) Service {
	return &service{r}
}

// Save a new purchase order to te database
func (s *service) Save(ctx context.Context, purchaseOrder domain.PurchaseOrder) (int, error) {
	// check if purchase order id already exists in the database
	existsPurchaseOrder := s.repo.ExistsPurchaseOrder(ctx, purchaseOrder.ID)
	if existsPurchaseOrder {
		return 0, ErrPurchaseOrderAlreadyExists
	}
	// check if buyers id exists
	existsBuyer := s.repo.ExistsBuyer(ctx, purchaseOrder.BuyerID)
	if !existsBuyer {
		return 0, ErrBuyerIDNotExists
	}
	// check if products record id exists
	existsProductsRecords := s.repo.ExistsProductsRecord(ctx, purchaseOrder.ProductRecordID)
	if !existsProductsRecords {
		return 0, ErrProductsRecordIDNotExits
	}

	// save purchase order into the database
	id, err := s.repo.Save(ctx, purchaseOrder)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// PurchaseOrdersByBuyer create a report with amount of purchase order by buyer
func (s *service) PurchaseOrdersByBuyer(ctx context.Context, buyerID int) ([]domain.PurchaseOrdersByBuyer, error) {

	// check if buyers id exists in the buyers database
	if buyerID != 0 {
		existsBuyer := s.repo.ExistsBuyer(ctx, buyerID)
		if !existsBuyer {
			return nil, ErrBuyerIDNotExists
		}
	}
	// process query to get purchase orders by buyer
	purchases, err := s.repo.PurchaseOrdersByBuyers(ctx, buyerID)
	if err != nil {
		return nil, err
	}
	// return purchases order report
	return purchases, nil

}
