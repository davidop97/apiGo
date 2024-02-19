package carries

import (
	"context"
	"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var (
	ErrIncorrectData = errors.New("incorrect data")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Carries, error)
	Save(ctx context.Context, w domain.Carries) (int, error)
	GetAllCarriesByLocality(ctx context.Context) ([]domain.LocalityCarries, error)
	GetAllCarriesByLocalityID(ctx context.Context, localityID int) (domain.LocalityCarries, error)
}

type service struct {
	rp Repository
}

func NewService(r Repository) Service {
	return &service{rp: r}
}

// GetAll is a method that returns all carries, returns empty list if there are no carries.
func (s *service) GetAll(ctx context.Context) ([]domain.Carries, error) {
	carriesList, err := s.rp.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return carriesList, nil
}

// Save is a method that saves a carry, returns error if the carry already exists or if the data is incorrect.
func (s *service) Save(ctx context.Context, c domain.Carries) (int, error) {
	for _, v := range []string{c.Address, c.Telephone, c.CID, c.CompanyName} {
		if v == "" {
			return 0, ErrIncorrectData
		}
	}
	if c.LocalityID < 0 {
		return 0, ErrIncorrectData
	}

	id, err := s.rp.Save(ctx, c)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetAllCarriesByLocality is a method that returns all carries by locality, returns empty list if there are no carries.
func (s *service) GetAllCarriesByLocality(ctx context.Context) ([]domain.LocalityCarries, error) {
	localityCarriesList, err := s.rp.GetAllCarriesByLocality(ctx)
	if err != nil {
		return nil, err
	}

	return localityCarriesList, nil
}

// GetAllCarriesByLocalityID is a method that returns all carries by locality ID, returns empty list if there are no carries.
func (s *service) GetAllCarriesByLocalityID(ctx context.Context, localityID int) (domain.LocalityCarries, error) {
	localityCarries, err := s.rp.GetAllCarriesByLocalityID(ctx, localityID)
	if err != nil {
		return localityCarries, err
	}

	return localityCarries, nil
}
