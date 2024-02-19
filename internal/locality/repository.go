package locality

import (
	"context"
	"database/sql"
	"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

var (
	ErrLocalityNotFound = errors.New("locality not found")
	ErrNoRows           = errors.New("no results for this request")
)

type Repository interface {
	GetLocality(ctx context.Context, id int) (domain.Locality, error)
	GetAll(ctx context.Context) ([]domain.Locality, error)
	Save(ctx context.Context, l domain.Locality) (int, error)
	Exists(ctx context.Context, cid int) bool
	GetReportSellers(ctx context.Context, id int) ([]domain.ReportSellers, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

// Get a locality using its id. Return an error if it doesn't exist.
func (r *repository) GetLocality(ctx context.Context, id int) (domain.Locality, error) {
	query := "SELECT id, postal_code, locality_name, province_name, country_name FROM locality WHERE id=?;"
	// Execute the query
	row := r.db.QueryRow(query, id)
	l := domain.Locality{}
	// Scan the row into the Locality struct
	err := row.Scan(&l.ID, &l.PostalCode, &l.LocalityName, &l.ProvinceName, &l.CountryName)
	if err != nil {
		//Check if the query returns a row or not
		if errors.Is(err, sql.ErrNoRows) {
			// Locality not found, return ErrNotFound
			return domain.Locality{}, ErrLocalityNotFound
		}
		// If occurs another error, return it to be controlled in the handler.
		return domain.Locality{}, err
	}
	// Everything is ok, return the requested Locality.
	return l, nil
}

// Get all the localities in the database. Return an error if it doesn't exist.
func (r *repository) GetAll(ctx context.Context) ([]domain.Locality, error) {
	query := "SELECT id, postal_code, locality_name, province_name, country_name FROM locality"
	//Execute the query
	rows, err := r.db.Query(query)
	//If an internal error occurs, it will be returned to be controlled in the handler.
	if err != nil {
		return nil, err
	}

	var localities []domain.Locality

	// Iterate over the rows, appending them to the list.
	for rows.Next() {
		l := domain.Locality{}
		// Scan the row into the locality struct.
		_ = rows.Scan(&l.ID, &l.PostalCode, &l.LocalityName, &l.ProvinceName, &l.CountryName)
		localities = append(localities, l)
	}
	//If another internal error occurs, it will be returned to be controlled in the handler.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(localities) == 0 {
		// If the list is empty, an ErrNotfound error will be returned to be controlled in the handler.
		return nil, ErrNoRows
	}
	// If everything is ok, the list of sellers will be returned.
	return localities, nil
}

// Save a new locality in the database. Return the id of the new locality and an error if it occurs.
func (r *repository) Save(ctx context.Context, l domain.Locality) (int, error) {
	query := "INSERT INTO locality (postal_code,locality_name, province_name, country_name) VALUES (?, ?, ?, ?)"
	// Prepare the query
	stmt, err := r.db.Prepare(query)
	// Check if an error occurs
	if err != nil {
		// If an error occurs, return 0 and the error to be controlled in the handler.
		return 0, err
	}

	// Execute the query
	res, err := stmt.Exec(l.PostalCode, l.LocalityName, l.ProvinceName, l.CountryName)
	if err != nil {
		// If an error occurs, return 0 and the error to be controlled in the handler.
		return 0, err
	}

	// Get the id of the new locality recently added.
	id, err := res.LastInsertId()
	if err != nil {
		// If an error occurs, return 0 and the error to be controlled in the handler.
		return 0, err
	}
	// Everything is ok, return the id of the new locality and nil as error.
	return int(id), nil
}

// Check if a Postal_code locality exists using its id. Return true if it exists and false if it doesn't.
func (r *repository) Exists(ctx context.Context, cid int) bool {
	query := "SELECT postal_code FROM locality WHERE postal_code=?;"
	row := r.db.QueryRow(query, cid)
	err := row.Scan(&cid)
	return err == nil
}

// Get a report of sellers in a locality using its id or in all localities if id is not provided. Return an error if it doesn't exist.
func (r *repository) GetReportSellers(ctx context.Context, id int) ([]domain.ReportSellers, error) {
	var args []interface{}
	query := `SELECT l.id AS locality_id, l.locality_name AS locality_name,l.postal_code, COUNT(s.id) AS seller_count
	FROM locality l
	LEFT JOIN sellers s ON l.id = s.locality_id `

	//Check if id is provided and if is greater than 0.
	if id > 0 {
		//If exists, add the id to the query and append it to the args slice.
		query += " WHERE l.id = ? "
		args = append(args, id)
	}
	//Group by locality id
	query += ` GROUP BY l.id;`

	// Execute the query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		//Check if the query returns a row or not.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRows
		}
		// If occurs another error, return it to be controlled in the handler.
		return nil, err
	}
	defer rows.Close()

	// Iterate the rows and append them to the slice of ReportSellers.
	var reportSellers []domain.ReportSellers
	for rows.Next() {
		var r domain.ReportSellers
		// Scan the row and append it to the slice. If an error occurs, return it to be controlled in the handler.
		if err := rows.Scan(&r.Locality_id, &r.Locality_name, &r.Postal_code, &r.Sellers_count); err != nil {
			return nil, err
		}
		//Generate the report
		reportSellers = append(reportSellers, r)
	}
	//Check if the query returns a row or not.
	if len(reportSellers) == 0 {
		// If the query doesn't return a row, return ErrNoRows.
		return nil, ErrNoRows
	}
	//If another errror occurs, return it to be controlled in the handler.
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Everything is ok, return the requested ReportSellers.
	return reportSellers, nil
}
