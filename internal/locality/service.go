package locality

import (
	"context"
	"errors"

	//"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var (
	ErrLocalityAlreadyExists = errors.New("locality already exists")
	ErrCannotSaveLocality    = errors.New("cannot add new locality")
)

type Service interface {
	GetLocalityByID(ctx context.Context, id int) (domain.Locality, error)
	GetAll(ctx context.Context) ([]domain.Locality, error)
	Save(ctx context.Context, l domain.Locality) (int, error)
	GetReportSellers(ctx context.Context, id int) ([]domain.ReportSellers, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

// Get a locality by id if exists.
func (s *service) GetLocalityByID(ctx context.Context, id int) (domain.Locality, error) {
	locality, err := s.r.GetLocality(ctx, id)
	return locality, err
}

// Get all the localities available.
func (s service) GetAll(ctx context.Context) ([]domain.Locality, error) {
	allLocalities, err := s.r.GetAll(ctx)
	return allLocalities, err
}

// Save a new locality or return an error if can not do that action.
func (s *service) Save(ctx context.Context, locality domain.Locality) (int, error) {
	//Check if the locality already exists.
	exist := s.r.Exists(ctx, locality.PostalCode)
	if exist {
		//If the locality already exists, return an error.
		return 0, ErrLocalityAlreadyExists
	}
	//Try to save the locality
	newLocality, err := s.r.Save(ctx, locality)
	if err != nil {
		//If an error occurs, return an error to be controlled in the handler.
		return 0, err
	}
	//If everything is ok, return the id of the new locality created.
	return newLocality, nil
}

// Get a report of sellers by locality.
func (s *service) GetReportSellers(ctx context.Context, id int) ([]domain.ReportSellers, error) {
	return s.r.GetReportSellers(ctx, id)
}
