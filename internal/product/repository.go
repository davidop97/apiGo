package product

import (
	"context"
	"database/sql"

	"github.com/davidop97/apiGo/internal/domain"
)

// Repository encapsulates the storage of a Product.
type Repository interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	Get(ctx context.Context, id int) (domain.Product, error)
	Exists(ctx context.Context, productCode string) bool
	Save(ctx context.Context, p domain.Product) (int, error)
	Update(ctx context.Context, p domain.Product) error
	Delete(ctx context.Context, id int) error
	CreateProductRecord(ctx context.Context, p domain.ProductRecordCreate) (int, error)
	GetProductRecord(ctx context.Context, idProduct int) ([]domain.ProductRecordGet, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Product, error) {
	query := "SELECT * FROM products;"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var products []domain.Product

	for rows.Next() {
		p := domain.Product{}
		_ = rows.Scan(&p.ID, &p.Description, &p.ExpirationRate, &p.FreezingRate, &p.Height, &p.Length, &p.Netweight, &p.ProductCode, &p.RecomFreezTemp, &p.Width, &p.ProductTypeID, &p.SellerID)
		products = append(products, p)
	}

	return products, nil
}

func (r *repository) Get(ctx context.Context, id int) (domain.Product, error) {
	query := "SELECT * FROM products WHERE id=?;"
	row := r.db.QueryRow(query, id)
	p := domain.Product{}
	err := row.Scan(&p.ID, &p.Description, &p.ExpirationRate, &p.FreezingRate, &p.Height, &p.Length, &p.Netweight, &p.ProductCode, &p.RecomFreezTemp, &p.Width, &p.ProductTypeID, &p.SellerID)
	if err != nil {
		return domain.Product{}, err
	}

	return p, nil
}

func (r *repository) Exists(ctx context.Context, productCode string) bool {
	query := "SELECT product_code FROM products WHERE product_code=?;"
	row := r.db.QueryRow(query, productCode)
	err := row.Scan(&productCode)
	return err == nil
}

func (r *repository) Save(ctx context.Context, p domain.Product) (int, error) {
	query := "INSERT INTO products(description,expiration_rate,freezing_rate,height,lenght,netweight,product_code,recommended_freezing_temperature,width,id_product_type,id_seller) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(p.Description, p.ExpirationRate, p.FreezingRate, p.Height, p.Length, p.Netweight, p.ProductCode, p.RecomFreezTemp, p.Width, p.ProductTypeID, p.SellerID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Update(ctx context.Context, p domain.Product) error {
	query := "UPDATE products SET description=?, expiration_rate=?, freezing_rate=?, height=?, lenght=?, netweight=?, product_code=?, recommended_freezing_temperature=?, width=?, id_product_type=?, id_seller=?  WHERE id=?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(p.Description, p.ExpirationRate, p.FreezingRate, p.Height, p.Length, p.Netweight, p.ProductCode, p.RecomFreezTemp, p.Width, p.ProductTypeID, p.SellerID, p.ID)
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
	query := "DELETE FROM products WHERE id=?"
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

// CreateProductRecord inserts a new product record into the database.
// It takes a context and a ProductRecordCreate struct as parameters.
// The ProductRecordCreate struct contains the details of the product record to be created.
// The function prepares an SQL statement for inserting the product record into the database.
// It then executes the SQL statement with the details from the ProductRecordCreate struct.
// If the SQL statement execution is successful, it retrieves the ID of the last inserted record.
// The function returns the ID of the created product record and an error if there is any.
func (r *repository) CreateProductRecord(ctx context.Context, p domain.ProductRecordCreate) (int, error) {
	query := "INSERT INTO productsRecord(last_update_date,purchase_price,sale_price,product_id) VALUES (STR_TO_DATE(?,'%Y-%m-%d'),?,?,?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(p.LastUpdate, p.PurchasePrice, p.SalePrice, p.ProductID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetProductRecord retrieves product records by product ID from the database.
// If idProduct is 0, it retrieves all product records.
// The function prepares an SQL statement for retrieving the product records from the database.
// It then executes the SQL statement with the product ID as a parameter.
// If the SQL statement execution is successful, it scans the result into a slice of ProductRecordGet structs.
// The function returns a slice of ProductRecordGet structs and an error if there is any.
func (r *repository) GetProductRecord(ctx context.Context, idProduct int) ([]domain.ProductRecordGet, error) {
	var query string
	var args []interface{}

	query = `
		SELECT
			p.id,
			p.description,
			COUNT(pr.id) AS records_count
		FROM products AS p
		LEFT JOIN productsRecord AS pr
			ON p.id = pr.product_id
	`
	// If idProduct is not 0, add a WHERE clause to the SQL statement
	if idProduct != 0 {
		query += " WHERE p.id = ?"
		args = append(args, idProduct)
	}

	// Add a GROUP BY clause to the SQL statement
	query += " GROUP BY p.id, p.description"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.ProductRecordGet
	for rows.Next() {
		var p domain.ProductRecordGet
		// Scan the result into a ProductRecordGet struct
		if err := rows.Scan(&p.ProductID, &p.Description, &p.RecordCount); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
