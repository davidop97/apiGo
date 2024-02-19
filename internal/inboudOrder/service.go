package inboudorder

import (
	"context"
	"errors"
	"log"

	"github.com/davidop97/apiGo/internal/domain"
	//"errors"
)

// Errors
var (
	ErrEmployeeDoesNotExists     = errors.New("Employee doesn`t exist")
	ErrInboundOrderAlreadyExists = errors.New("Inboud order already exist")
	ErrWarehouseDoesNotExists    = errors.New("Warehouse doesn`t exist")
)

type Service interface {
	GetAllReports(ctx context.Context) ([]Report, error)
	GenerateReport(ctx context.Context, employeeID int) (Report, error)
	CreateInboundOrder(ctx context.Context, order domain.InboudOrder) (id int, err error)
}

// Struct contains repository
type service struct {
	repo Repository
}

// Constructor from service struct, receive repository
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetAllReports(ctx context.Context) ([]Report, error) {
	return s.repo.GetAllReports(ctx)
}

func (s *service) GenerateReport(ctx context.Context, employeeID int) (Report, error) {
	return s.repo.GenerateReport(ctx, employeeID)
}

func (s *service) CreateInboundOrder(ctx context.Context, order domain.InboudOrder) (id int, err error) {
	// Verificar si el Employee existe
	if !s.repo.ExistsEmployee(ctx, order.EmployeeID) {
		// Employee doesn`t exists
		err = ErrEmployeeDoesNotExists
		return
	}

	// Verificar si warehouse existe
	if !s.repo.ExistsWarehouse(ctx, order.WarehouseID) {
		// Warehouse doesn`t exists
		err = ErrWarehouseDoesNotExists
		return
	}

	// Verificar si la Inbound Order ya existe
	if s.repo.ExistsInboundOrder(ctx, order.OrderNumber) {
		// Employee already exists
		err = ErrInboundOrderAlreadyExists
		return
	}

	// Crear la Inbound Order
	id, err = s.repo.Save(ctx, order)
	if err != nil {
		log.Printf("Error saving Inbound Order: %v", err) // Imprimir el error en el log
		return
	}
	return
}
