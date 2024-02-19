package employee

import (
	"context"
	"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var ErrEmployeeAlreadyExists = errors.New("Employee already exist")

// Interface Service with methods
type Service interface {
	GetAllEmployees(ctx context.Context) ([]domain.Employee, error)
	GetEmployeeByID(ctx context.Context, id int) (domain.Employee, error)
	SaveEmployee(ctx context.Context, employee domain.Employee) (int, error)
	UpdateEmployee(ctx context.Context, employee domain.Employee) error
	DeleteEmployee(ctx context.Context, id int) error
}

// Struct contains repository
type service struct {
	repo Repository
}

// Constructor from service struct, receive repository
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Function get all employees
func (s *service) GetAllEmployees(ctx context.Context) (employees []domain.Employee, err error) {
	employees, err = s.repo.GetAll(ctx)

	return
}

// Function get one employee
func (s *service) GetEmployeeByID(ctx context.Context, id int) (employee domain.Employee, err error) {
	employee, err = s.repo.Get(ctx, id)
	return
}

// Function save employee
func (s *service) SaveEmployee(ctx context.Context, employee domain.Employee) (id int, err error) {
	// Verify if employee exists
	if s.repo.Exists(ctx, employee.CardNumberID) {
		// Employee already exists
		err = ErrEmployeeAlreadyExists
		return
	}
	// Insertion Employee
	id, err = s.repo.Save(ctx, employee)
	return
}

// Function update employee
func (s *service) UpdateEmployee(ctx context.Context, employee domain.Employee) (err error) {
	// Get actual employee
	currentEmployee, err := s.repo.Get(ctx, employee.ID)
	if err != nil {
		return
	}

	// - If emplooye number is being updated, check if the new one already exists
	if employee.CardNumberID != currentEmployee.CardNumberID {
		exists := s.repo.Exists(ctx, employee.CardNumberID)
		if exists {
			err = ErrEmployeeAlreadyExists
			return
		}
	}

	err = s.repo.Update(ctx, employee)

	return
}

func (s *service) DeleteEmployee(ctx context.Context, id int) (err error) {
	err = s.repo.Delete(ctx, id)

	return
}
