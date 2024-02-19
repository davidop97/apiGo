package purchase_order

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/davidop97/apiGo/internal/domain"
)

// Repository is an interface that defines the methods for a repository
// that handles purchase order data.
type Repository interface {
	// Save saves a purchase order to the database.
	Save(ctx context.Context, po domain.PurchaseOrder) (int, error)
	// ExistsPurchaseOrder checks if a purchase order already exists in the database.
	ExistsPurchaseOrder(ctx context.Context, purchaseOrder int) bool
	// ExistsBuyer checks if a buyer already exists in the database.
	ExistsBuyer(ctx context.Context, id int) bool
	// ExistsProductsRecord checks if a record of products already exists in the database.
	ExistsProductsRecord(ctx context.Context, id int) bool
	// PurchaseOrdersByBuyers returns all purchase orders made by a specific buyer.
	PurchaseOrdersByBuyers(ctx context.Context, buyerID int) ([]domain.PurchaseOrdersByBuyer, error)
}

// repository is the concrete implementation of the Repository interface.
type repository struct {
	db *sql.DB
}

// NewRepository creates a new instance of the repository.
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// check if id field of purchase_orders table exists into the database
func (r *repository) ExistsPurchaseOrder(ctx context.Context, purchaseOrderID int) bool {
	// query to select the id field from the purchase_orders table where the id matches the input id
	query := "SELECT id FROM purchase_orders WHERE id=?;"
	row := r.db.QueryRow(query, purchaseOrderID)
	// scan the result of the query into the input id variable
	err := row.Scan(&purchaseOrderID)
	// return true if the scan was successful (i.e., if the id exists in the table), and false otherwise
	return err == nil
}

// check if id field of buyers table exists into the database
func (r *repository) ExistsBuyer(ctx context.Context, id int) bool {
	// SQL query to select the id field from the buyers table where the id matches the input id
	query := "SELECT id FROM buyers WHERE id=?;"
	row := r.db.QueryRow(query, id)
	// scan the result of the query into the input id variable
	err := row.Scan(&id)
	// return true if the scan was successful (i.e., if the id exists in the table), and false otherwise
	return err == nil
}

// check if id of productsRecord table exists into the database
func (r *repository) ExistsProductsRecord(ctx context.Context, id int) bool {
	// queryery to select the id field from the productsRecord table where the id matches the input id
	query := "SELECT id FROM productsRecord WHERE id=?;"
	row := r.db.QueryRow(query, id)
	// scan the result of the query into the input id variable
	err := row.Scan(&id)
	// return true if the scan was successful (i.e., if the id exists in the table), and false otherwise
	return err == nil
}

// Create purchase order and save it into the database
func (r *repository) Save(ctx context.Context, po domain.PurchaseOrder) (int, error) {
	query := "INSERT INTO purchase_orders(order_number, order_date, tracking_code, buyer_id, product_record_id, order_status_id) VALUES (?, STR_TO_DATE(?,'%Y-%m-%d'), ?, ?, ?, ?)"
	// prepare query
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	// execute the query
	res, err := stmt.Exec(&po.OrderNumber, &po.OrderDate, &po.TrackingCode, &po.BuyerID, &po.ProductRecordID, &po.OrderStatusID)
	if err != nil {
		return 0, err
	}
	// id of the new purchase order added
	id, err := res.LastInsertId()
	if err != nil {
		return 0, nil
	}
	// return id and nil of the object created
	return int(id), nil
}

// Report numbers of purchase orders per buyer
func (r *repository) PurchaseOrdersByBuyers(ctx context.Context, buyerID int) ([]domain.PurchaseOrdersByBuyer, error) {
	// prepare query
	query := `SELECT
			  	b.id,
				b.card_number_id,
				b.first_name,
				b.last_name,
				COUNT(po.id) AS purchase_orders_count
			  FROM
			   buyers b
			  LEFT JOIN
			   Purchase_Orders po ON b.id = po.buyer_id
			  `
	// if buyer id is present, filter with that id
	if buyerID != 0 {
		query += fmt.Sprintf(" WHERE b.id = %d", buyerID)
	}
	// complete the rest of the query group by the fields require
	query += " GROUP BY b.id, b.card_number_id, b.first_name, b.last_name"

	// Execute query
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []domain.PurchaseOrdersByBuyer

	// Iterate over results
	for rows.Next() {
		// Create a new PurchaseOrdersByBuyer struct to store the data for each row
		po := domain.PurchaseOrdersByBuyer{}

		// Scan the column values of the current row into the PurchaseOrdersByBuyer struct
		err := rows.Scan(&po.ID, &po.CardNumberID, &po.FirstName, &po.LastName, &po.PurchaseOrdersCount)
		if err != nil {
			return nil, err
		}
		// Append the PurchaseOrdersByBuyer struct to the results slice
		results = append(results, po)

	}
	// Check if any errors occurred during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return the results slice containing the purchase order data
	return results, nil
}
