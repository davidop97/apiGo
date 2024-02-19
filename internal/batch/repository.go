package batch

import (
	"context"
	"database/sql"
	"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var (
	ErrProductNotFound = errors.New("associated product not found")
	ErrSectionNotFound = errors.New("associated section not found")
)

type Repository interface {
	GetAll(ctx context.Context) ([]domain.ProductBatch, error)
	Save(ctx context.Context, b domain.ProductBatch) (int, error)
	Exists(ctx context.Context, batchNumber int) bool
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

// GetAll returns all Product Batches stored in the database
func (r *repository) GetAll(ctx context.Context) (batches []domain.ProductBatch, err error) {
	query := "SELECT id, batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id FROM productBatches;"
	rows, err := r.db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		b := domain.ProductBatch{}
		_ = rows.Scan(&b.ID, &b.BatchNumber, &b.CurrentQuantity, &b.CurrentTemperature, &b.DueDate, &b.InitialQuantity, &b.ManufacturingDate, &b.ManufacturingHour, &b.MinimumTemperature, &b.ProductID, &b.SectionID)
		batches = append(batches, b)
	}
	return
}

// Save stores a new Product Batch in the database
func (r *repository) Save(ctx context.Context, b domain.ProductBatch) (int, error) {
	// Check if foreign keys exist
	// - check if associated product exists
	exists := r.productExists(ctx, b.ProductID)
	if !exists {
		return 0, ErrProductNotFound
	}
	// - check if associated section exists
	exists = r.sectionExists(ctx, b.SectionID)
	if !exists {
		return 0, ErrSectionNotFound
	}

	// Prepare query
	query := "INSERT INTO productBatches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	// Run query
	res, err := stmt.Exec(&b.BatchNumber, &b.CurrentQuantity, &b.CurrentTemperature, &b.DueDate, &b.InitialQuantity, &b.ManufacturingDate, &b.ManufacturingHour, &b.MinimumTemperature, &b.ProductID, &b.SectionID)
	if err != nil {
		return 0, err
	}

	// Return result
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Exists checks wether a given batch number is already stored in the database
func (r *repository) Exists(ctx context.Context, batchNumber int) bool {
	query := "SELECT batch_number FROM productBatches WHERE batch_number=?;"
	row := r.db.QueryRow(query, batchNumber)
	err := row.Scan(&batchNumber)
	return err == nil
}

// productExists is an auxiliary function that checks if a product id exists in the database
func (r *repository) productExists(ctx context.Context, id int) bool {
	query := "SELECT id FROM products WHERE id=?;"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&id)
	return err == nil
}

// sectionExists is an auxiliary function that checks if a section id exists in the database
func (r *repository) sectionExists(ctx context.Context, id int) bool {
	query := "SELECT id FROM sections WHERE id=?;"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&id)
	return err == nil
}
