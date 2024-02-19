package seller

import (
	"context"
	"database/sql"
	"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var (
	ErrNotFound = errors.New("seller not found")
)

// Repository encapsulates the storage of a Seller.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Seller, error)
	Get(ctx context.Context, id int) (domain.Seller, error)
	Exists(ctx context.Context, cid int) bool
	Save(ctx context.Context, s domain.Seller) (int, error)
	Update(ctx context.Context, s domain.Seller) error
	Delete(ctx context.Context, id int) error
	GetLocalityIdFromSeller(ctx context.Context, id int) bool
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

// Get all the sellers in the database. Return an error if the list is empty
// or another internal error occurs, it will be returned to be controlled in the handler.
func (r *repository) GetAll(ctx context.Context) ([]domain.Seller, error) {
	query := "SELECT id, cid, company_name, address, telephone, locality_id FROM sellers"
	//Execute the query.
	rows, err := r.db.Query(query)
	//If an internal error occurs, it will be returned to be controlled in the handler.
	if err != nil {
		return nil, err
	}

	var sellers []domain.Seller
	// Iterate over the rows.
	for rows.Next() {
		s := domain.Seller{}
		// Scan the row into the Seller struct.
		_ = rows.Scan(&s.ID, &s.CID, &s.CompanyName, &s.Address, &s.Telephone, &s.IDLocality)
		sellers = append(sellers, s)
	}
	//If another internal error occurs, it will be returned to be controlled in the handler.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(sellers) == 0 {
		// If the list is empty, an ErrNotfound error will be returned to be controlled in the handler.
		return nil, ErrNotFound
	}
	// If everything is ok, the list of sellers will be returned.
	return sellers, nil
}

// Get a seller using its id. Return an error if it doesn't exist.
func (r *repository) Get(ctx context.Context, id int) (domain.Seller, error) {
	query := "SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id=?;"
	// Execute the query.
	row := r.db.QueryRow(query, id)
	s := domain.Seller{}
	// Scan the row into the Seller struct.
	err := row.Scan(&s.ID, &s.CID, &s.CompanyName, &s.Address, &s.Telephone, &s.IDLocality)
	if err != nil {
		//Check if the query returns a row or not.
		if errors.Is(err, sql.ErrNoRows) {
			// Seller not found, return ErrNotFound.
			return domain.Seller{}, ErrNotFound
		}
		// If occurs another error, return it to be controlled in the handler.
		return domain.Seller{}, err
	}
	// Everything is ok, return the requested seller.
	return s, nil
}

func (r *repository) Exists(ctx context.Context, cid int) bool {
	query := "SELECT cid FROM sellers WHERE cid=?;"
	row := r.db.QueryRow(query, cid)
	err := row.Scan(&cid)
	return err == nil
}

// Save a seller in the database. Return the last inserted id or an error if it occurs
// to be controlled in the handler.
func (r *repository) Save(ctx context.Context, s domain.Seller) (int, error) {
	query := "INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)"
	// Prepare the query.
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	// Execute the query.
	res, err := stmt.Exec(s.CID, s.CompanyName, s.Address, s.Telephone, s.IDLocality)
	if err != nil {
		return 0, err
	}
	// Get the last inserted id.
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	// Everything is ok, return the last inserted id.
	return int(id), nil
}

// Update a seller in the database. Return an error if it doesn't exist or another internal error occurs
// to be controlled in the handler.
func (r *repository) Update(ctx context.Context, s domain.Seller) error {
	query := "UPDATE sellers SET cid=?, company_name=?, address=?, telephone=?, locality_id=? WHERE id=?"
	// Prepare the query.
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	// Execute the query.
	res, err := stmt.Exec(s.CID, s.CompanyName, s.Address, s.Telephone, s.ID, s.IDLocality)
	if err != nil {
		return err
	}

	// Check if the seller was updated.
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	// Everything is ok, return nil.
	return nil
}

// Delete a seller from the database. Return an error if it doesn't exist or another internal error occurs
// to be controlled in the handler.
func (r *repository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM sellers WHERE id=?"
	// Prepare the query.
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect < 1 {
		return ErrNotFound
	}

	return nil
}

// Check if a locality_id exists using its id. Return true if it exists and false otherwise.
func (r *repository) GetLocalityIdFromSeller(ctx context.Context, id int) bool {
	query := "SELECT id FROM locality WHERE id = ?;"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&id)
	return err == nil
}
