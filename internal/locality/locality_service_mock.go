package locality

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"

	"github.com/stretchr/testify/mock"
)

// ServiceMock struct. Mock of the service.
type ServiceMock struct {
	mock.Mock
}

// NewMockService function. Create a new mock of the service.
func NewMockService() *ServiceMock {
	return &ServiceMock{}
}

// GetLocalityByID function. Mock of the GetLocalityByID function. Get a locality by id if exists or an error.
func (s *ServiceMock) GetLocalityByID(ctx context.Context, id int) (i domain.Locality, err error) {
	//Get the locality by id.
	args := s.Called(ctx, id)
	//Return the locality and the error if exists.
	return args.Get(0).(domain.Locality), args.Error(1)
}

// GetAllLocalities function. Mock of the GetAll function. Get all the localities available or an error.
func (s *ServiceMock) GetAll(ctx context.Context) (l []domain.Locality, err error) {
	//Get all the localities.
	args := s.Called(ctx)
	//Return an array of localities and the error if exists.
	return args.Get(0).([]domain.Locality), args.Error(1)
}

// Save function. Mock of the Save function. Save a new locality or return an error if can not do that action.
func (s *ServiceMock) Save(ctx context.Context, locality domain.Locality) (id int, err error) {
	//Try to Save the locality.
	args := s.Called(ctx, locality)
	//Return the id of the locality if could be added and the error if exists.
	return args.Int(0), args.Error(1)
}

// GetReportSellers function. Mock of the GetReportSellers function. Get a report of sellers by locality or an error.
func (s *ServiceMock) GetReportSellers(ctx context.Context, id int) (r []domain.ReportSellers, err error) {
	//Get the report of sellers by locality.
	args := s.Called(ctx, id)
	//Return the report of sellers and the error if exists.
	return args.Get(0).([]domain.ReportSellers), args.Error(1)
}
