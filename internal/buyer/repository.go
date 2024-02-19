package buyer

import (
	"context"
	"database/sql"

	"github.com/davidop97/apiGo/internal/domain"
)

// Repository encapsulates the storage of a buyer.
type Repository interface {
	// GetAll obtains all buyers
	GetAll(ctx context.Context) ([]domain.Buyer, error)
	// Get returns a single buyer by its ID. If it doesn't exist
	Get(ctx context.Context, id int) (domain.Buyer, error)
	// Exists checks if buyer with certain card number id exists
	Exists(ctx context.Context, cardNumberID string) bool
	// Save adds a new buyer to the repository
	Save(ctx context.Context, b domain.Buyer) (int, error)
	// Update a buyer
	Update(ctx context.Context, b domain.Buyer) error
	// Delete deletes a buyer by id
	Delete(ctx context.Context, id int) error
}

// repository is the concrete implementation of the Repository interface.
type repository struct {
	db *sql.DB
}

// NewRepository creates a new instance of the repository
func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

// GettAll obtains all buyers
func (r *repository) GetAll(ctx context.Context) ([]domain.Buyer, error) {
	// query to select all buyers
	query := "SELECT * FROM buyers"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var buyers []domain.Buyer
	// loop rows
	for rows.Next() {
		b := domain.Buyer{}
		// scan the result of the query
		_ = rows.Scan(&b.ID, &b.CardNumberID, &b.FirstName, &b.LastName)
		// collect al the buyers
		buyers = append(buyers, b)
	}

	// return slice of buyers
	return buyers, nil
}

// Get gets a single buyer by its ID
func (r *repository) Get(ctx context.Context, id int) (domain.Buyer, error) {
	// query to get a buyer by ide
	query := "SELECT * FROM buyers WHERE id = ?;"
	row := r.db.QueryRow(query, id)
	b := domain.Buyer{}
	// scan the result of the query
	err := row.Scan(&b.ID, &b.CardNumberID, &b.FirstName, &b.LastName)
	if err != nil {
		return domain.Buyer{}, err
	}

	// returns a buyer if not error
	return b, nil
}

// Exists check if buyer with certain card number id exists
func (r *repository) Exists(ctx context.Context, cardNumberID string) bool {
	// query to obtain a buyer by id
	query := "SELECT card_number_id FROM buyers WHERE card_number_id=?;"
	row := r.db.QueryRow(query, cardNumberID)
	// scan result of the query
	err := row.Scan(&cardNumberID)
	// if any error occurs or there is no record with that id it returns false
	return err == nil
}

// Save a new buyer into the database
func (r *repository) Save(ctx context.Context, b domain.Buyer) (int, error) {
	// query to insert a new buyer
	query := "INSERT INTO buyers(card_number_id,first_name,last_name) VALUES (?,?,?)"
	// prepare the query
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	// Exec executes a prepared statement with the given arguments.
	res, err := stmt.Exec(&b.CardNumberID, &b.FirstName, &b.LastName)
	if err != nil {
		return 0, err
	}

	// get the last id inserted into the database
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// returns the present id
	return int(id), nil
}

// Update an existing buyer in the database
func (r *repository) Update(ctx context.Context, b domain.Buyer) error {
	query := "UPDATE buyers SET first_name=?, last_name=?  WHERE id=?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	// Exec executes a prepared statement with the given arguments
	res, err := stmt.Exec(&b.FirstName, &b.LastName, &b.ID)
	if err != nil {
		return err
	}

	// RowsAffected returns the number of rows affected by the update
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	// returns nil if there is no error
	return nil
}

// Delete removes a buyer from the database
func (r *repository) Delete(ctx context.Context, id int) error {
	// query searching a buyer by id
	query := "DELETE FROM buyers WHERE id = ?"
	// prepare query
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	// Exec executes a prepared statement with the given arguments
	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	// RowsAffected returns the number of rows affected by the delete.
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// check if there any rows affect
	if affect < 1 {
		return ErrNotFound
	}

	return nil
}
