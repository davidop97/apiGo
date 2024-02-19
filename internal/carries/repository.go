package carries

import (
	"context"
	"database/sql"
	"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var (
	ErrDuplicateCarry          = errors.New("warehouse already exists")
	ErrLocalityCarriesNotFound = errors.New("locality carries not found")
)

type Repository interface {
	GetAll(ctx context.Context) ([]domain.Carries, error)
	Save(ctx context.Context, w domain.Carries) (int, error)
	GetAllCarriesByLocality(ctx context.Context) ([]domain.LocalityCarries, error)
	GetAllCarriesByLocalityID(ctx context.Context, localityID int) (domain.LocalityCarries, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

// GetAll is a method that returns all carries, returns empty list if there are no carries.
func (r *repository) GetAll(ctx context.Context) ([]domain.Carries, error) {
	query := "SELECT * FROM carries"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var carriesList []domain.Carries

	for rows.Next() {
		c := domain.Carries{}
		err = rows.Scan(&c.ID, &c.CID, &c.CompanyName, &c.Address, &c.Telephone, &c.LocalityID)
		if err != nil {
			return nil, err
		}
		carriesList = append(carriesList, c)
	}

	defer rows.Close()

	return carriesList, nil
}

// Save is a method that saves a carry, returns error if the carry already exists or if the data is incorrect.
func (r *repository) Save(ctx context.Context, c domain.Carries) (int, error) {
	if r.Exists(ctx, c.CID) {
		return 0, ErrDuplicateCarry
	} else if !r.LocalityExists(ctx, c.LocalityID) {
		return 0, ErrLocalityCarriesNotFound
	}

	query := "INSERT INTO carries (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?);"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(c.CID, c.CompanyName, c.Address, c.Telephone, c.LocalityID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Exists is a method that returns true if the carry exists, false otherwise.
func (r *repository) Exists(ctx context.Context, cid string) bool {
	query := "SELECT cid FROM carries WHERE cid=?;"
	row := r.db.QueryRow(query, cid)
	err := row.Scan(&cid)
	return err == nil
}

// Exists is a method that returns true if the carry exists, false otherwise.
func (r *repository) LocalityExists(ctx context.Context, id int) bool {
	query := "SELECT postal_code FROM locality WHERE postal_code=?;"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&id)
	return err == nil
}

// GetAllCarriesByLocality is a method that returns all carries by locality, returns empty list if there are no carries.
func (r *repository) GetAllCarriesByLocality(ctx context.Context) ([]domain.LocalityCarries, error) {
	query := `SELECT localities.postal_code, localities.locality_name, COUNT(*) FROM melisprint.carries as carries
			JOIN melisprint.locality as localities ON carries.locality_id = localities.postal_code
			GROUP BY localities.postal_code, localities.locality_name;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var carriesList []domain.LocalityCarries

	for rows.Next() {
		c := domain.LocalityCarries{}
		err = rows.Scan(&c.LocalityID, &c.LocalityName, &c.CarriesCount)
		if err != nil {
			return nil, err
		}
		carriesList = append(carriesList, c)
	}

	defer rows.Close()

	return carriesList, nil
}

// GetAllCarriesByLocalityID is a method that returns all carries by locality, returns empty list if there are no carries.
func (r *repository) GetAllCarriesByLocalityID(ctx context.Context, localityID int) (domain.LocalityCarries, error) {
	var lc domain.LocalityCarries

	if !r.LocalityExists(ctx, localityID) {
		return lc, ErrLocalityCarriesNotFound
	}

	query := `SELECT localities.postal_code, localities.locality_name, COUNT(*) FROM melisprint.carries as carries
			JOIN melisprint.locality as localities ON carries.locality_id = localities.postal_code
			WHERE postal_code = ?
			GROUP BY localities.postal_code, localities.locality_name;`

	row := r.db.QueryRow(query, localityID)

	err := row.Scan(&lc.LocalityID, &lc.LocalityName, &lc.CarriesCount)
	if err != nil {
		return lc, err
	}

	return lc, nil
}
