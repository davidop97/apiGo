package locality

import (
	"context"

	"github.com/davidop97/apiGo/internal/domain"

	"github.com/stretchr/testify/mock"
)

// RepositoryMock struct. Mock of the repository.
type RepositoryMock struct {
	mock.Mock
}

// NewMockRepository function. Create a new mock of the repository.
func NewMockRepository() *RepositoryMock {
	return &RepositoryMock{}
}

// GetLocality function. Mock of the GetLocality function. Get a locality by id if exists or an error.
func (r *RepositoryMock) GetLocality(ctx context.Context, id int) (i domain.Locality, err error) {
	//Get the locality by id.
	args := r.Called(ctx, id)
	//Return the locality and the error if exists.
	return args.Get(0).(domain.Locality), args.Error(1)
}

// GetAll function. Mock of the GetAll function. Get all the localities available or an error.
func (r *RepositoryMock) GetAll(ctx context.Context) (l []domain.Locality, err error) {
	//Get all the localities.
	args := r.Called(ctx)
	//Return an array of localities and the error if exists.
	return args.Get(0).([]domain.Locality), args.Error(1)
}

// Save function. Mock of the Save function. Save a new locality or return an error if can not do that action.
func (r *RepositoryMock) Save(ctx context.Context, l domain.Locality) (id int, err error) {
	//Try to Save the locality.
	args := r.Called(ctx, l)
	//Return the id of the locality if could be added and the error if exists.
	return args.Int(0), args.Error(1)
}

// Exists function. Mock of the Exists function. Check if a locality exists or not.
func (r *RepositoryMock) Exists(ctx context.Context, cid int) bool {
	//Check if the locality exists.
	args := r.Called(ctx, cid)
	//Return true if the locality exists. Otherwise returns false.
	return args.Bool(0)
}

// GetReportSellers function. Mock of the GetReportSellers function. Get a report of sellers by locality id if exists or an error.
func (r *RepositoryMock) GetReportSellers(ctx context.Context, id int) (l []domain.ReportSellers, err error) {
	//Get the report of sellers by locality id.
	args := r.Called(ctx, id)
	//Return the report of sellers and the error if exists.
	return args.Get(0).([]domain.ReportSellers), args.Error(1)
}
