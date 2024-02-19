package inboudorder

import (
	"context"
	"database/sql"
	"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

var ErrEmployeeNotFound = errors.New("section not found")

type Repository interface {
	GetAllReports(ctx context.Context) ([]Report, error)
	Exists(ctx context.Context, employeeID int) (*domain.Employee, error)
	GenerateReport(ctx context.Context, employeeID int) (report Report, err error)
	ExistsEmployee(ctx context.Context, employeeID int) bool
	ExistsInboundOrder(ctx context.Context, orderNumber string) bool
	ExistsWarehouse(ctx context.Context, warehouseID int) bool
	Save(ctx context.Context, i domain.InboudOrder) (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

type Report struct {
	*domain.Employee
	InboudOrdersCount int `json:"inboud_orders_count"`
}

func (r *repository) Exists(ctx context.Context, employeeID int) (*domain.Employee, error) {
	query := "SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees WHERE id = ?"
	var employee domain.Employee
	err := r.db.QueryRowContext(ctx, query, employeeID).Scan(
		&employee.ID,
		&employee.CardNumberID,
		&employee.FirstName,
		&employee.LastName,
		&employee.WarehouseID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No se encontró el empleado
			return nil, nil
		}
		// Ocurrió un error diferente
		return nil, err
	}

	return &employee, nil
}

func (r *repository) GenerateReport(ctx context.Context, employeeID int) (report Report, err error) {
	// Escenario 1 y 2: Validar si el Employee existe
	employee, err := r.Exists(ctx, employeeID)
	if err != nil {
		return
	}

	if employee == nil {
		// Escenario 1: Employee no existe
		err = ErrEmployeeNotFound
		return
	}

	// Escenario 2: Employee sí existe
	// Obtener el conteo de órdenes de entrada asociadas al empleado
	query := "SELECT COUNT(*) FROM inboudOrders WHERE employee_id = ?"
	err = r.db.QueryRowContext(ctx, query, employeeID).Scan(&report.InboudOrdersCount)
	if err != nil {
		return
	}

	// Configurar los campos del reporte con la información de employee
	report.Employee = employee

	return
}

func (r *repository) GetAllReports(ctx context.Context) ([]Report, error) {
	// Obtener todos los empleados
	query := "SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []Report

	// Iterar sobre los empleados y generar informe para cada uno
	for rows.Next() {
		var employee domain.Employee
		err := rows.Scan(
			&employee.ID,
			&employee.CardNumberID,
			&employee.FirstName,
			&employee.LastName,
			&employee.WarehouseID,
		)
		if err != nil {
			return nil, err
		}

		// Generar informe para el empleado actual
		report := Report{Employee: &employee}

		// Obtener el conteo de órdenes de entrada asociadas al empleado
		orderCountQuery := "SELECT COUNT(*) FROM inboudOrders WHERE employee_id = ?"
		err = r.db.QueryRowContext(ctx, orderCountQuery, employee.ID).Scan(&report.InboudOrdersCount)
		if err != nil {
			return nil, err
		}

		// Agregar el informe a la lista de informes
		reports = append(reports, report)
	}

	return reports, nil
}

func (r *repository) ExistsEmployee(ctx context.Context, employeeID int) bool {
	query := "SELECT id FROM employees WHERE id=?;"
	row := r.db.QueryRow(query, employeeID)
	err := row.Scan(&employeeID)
	return err == nil
}

func (r *repository) ExistsInboundOrder(ctx context.Context, orderNumber string) bool {
	query := "SELECT order_number FROM inboudOrders WHERE order_number=?;"
	row := r.db.QueryRow(query, orderNumber)
	err := row.Scan(&orderNumber)
	return err == nil
}

func (r *repository) ExistsWarehouse(ctx context.Context, warehouseID int) bool {
	query := "SELECT id FROM warehouses WHERE id=?;"
	row := r.db.QueryRow(query, warehouseID)
	err := row.Scan(&warehouseID)
	return err == nil
}

func (r *repository) Save(ctx context.Context, i domain.InboudOrder) (int, error) {
	query := "INSERT INTO inboudOrders(order_date, order_number, employee_id, product_batch_id, warehouse_id) VALUES (?,?,?,?,?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(&i.OrderDate, &i.OrderNumber, &i.EmployeeID, &i.ProductBatchID, &i.WarehouseID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
