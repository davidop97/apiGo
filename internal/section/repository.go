package section

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var (
	ErrNotFound = errors.New("section not found")
)

type ProdCountResponse struct {
	ID            int `json:"id"`
	SectionNumber int `json:"section_number"`
	ProductCount  int `json:"product_count"`
}

// Repository encapsulates the storage of a section.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Section, error)
	Get(ctx context.Context, id int) (domain.Section, error)
	Exists(ctx context.Context, cid int) bool
	Save(ctx context.Context, s domain.Section) (int, error)
	Update(ctx context.Context, s domain.Section) error
	Delete(ctx context.Context, id int) error
	ProductCount(ctx context.Context, id int) ([]ProdCountResponse, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Section, error) {
	query := "SELECT id, section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, product_type_id FROM sections;"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var sections []domain.Section

	for rows.Next() {
		s := domain.Section{}
		_ = rows.Scan(&s.ID, &s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
		sections = append(sections, s)
	}

	return sections, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Section, error) {
	query := "SELECT id, section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, product_type_id FROM sections WHERE id=?;"
	row := r.db.QueryRow(query, id)
	s := domain.Section{}
	err := row.Scan(&s.ID, &s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return domain.Section{}, ErrNotFound
		default:
			return domain.Section{}, err
		}
	}
	return s, nil
}

func (r *repository) Exists(ctx context.Context, sectionNumber int) bool {
	query := "SELECT section_number FROM sections WHERE section_number=?;"
	row := r.db.QueryRow(query, sectionNumber)
	err := row.Scan(&sectionNumber)
	return err == nil
}

func (r *repository) Save(ctx context.Context, s domain.Section) (int, error) {
	query := "INSERT INTO sections (section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, id_product_type) VALUES (?, ?, ?, ?, ?, ?, ?, ?);"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(&s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Update(ctx context.Context, s domain.Section) error {
	query := "UPDATE sections SET section_number=?, current_temperature=?, minimum_temperature=?, current_capacity=?, minimum_capacity=?, maximum_capacity=?, warehouse_id=?, id_product_type=? WHERE id=?;"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(&s.SectionNumber, &s.CurrentTemperature, &s.MinimumTemperature, &s.CurrentCapacity, &s.MinimumCapacity, &s.MaximumCapacity, &s.WarehouseID, &s.ProductTypeID, &s.ID)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	// Delete associated ProductBatches
	err := r.deleteBatches(ctx, id)
	if err != nil {
		return err
	}

	// Delete section
	query := "DELETE FROM sections WHERE id=?;"
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

// deleteBatches deletes ProductBatches associated with a product id
func (r *repository) deleteBatches(ctx context.Context, id int) (err error) {
	query := "DELETE FROM productBatches WHERE section_id=?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return
	}
	_, err = stmt.Exec(id)
	return
}

// ProductCount returns the number of products contained in each section
func (r *repository) ProductCount(ctx context.Context, id int) (l []ProdCountResponse, err error) {
	// Build query
	// - create id filter sql expression if id is given
	idFilter := ""
	if id != 0 {
		idFilter = fmt.Sprintf("WHERE sections.id=%d ", id)
	}
	// - create query string
	query := "SELECT sections.id, sections.section_number, sum(productBatches.current_quantity) as products_count " +
		"FROM sections LEFT JOIN productBatches " +
		"ON sections.id = productBatches.section_id " +
		idFilter +
		"GROUP BY sections.id;"

	// Run query
	rows, err := r.db.Query(query)
	if err != nil {
		return
	}

	// Process response
	for rows.Next() {
		res := ProdCountResponse{}
		prodCount := sql.NullInt64{}
		_ = rows.Scan(&res.ID, &res.SectionNumber, &prodCount)
		if prodCount.Valid {
			res.ProductCount = int(prodCount.Int64)
		} else {
			res.ProductCount = 0
		}
		l = append(l, res)
	}
	if id != 0 && len(l) < 1 {
		err = ErrNotFound
	}
	return
}
