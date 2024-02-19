package section

import (
	"context"
	"errors"

	"github.com/davidop97/apiGo/internal/domain"
)

// Errors
var (
	ErrDuplicateSectNumber = errors.New("duplicate section number")
)

type Service interface {
	GetAll(ctx context.Context) ([]domain.Section, error)
	Get(ctx context.Context, id int) (domain.Section, error)
	Save(ctx context.Context, sect domain.Section) (id int, err error)
	Delete(ctx context.Context, id int) (err error)
	Update(ctx context.Context, sect domain.Section) (err error)
	ProductCount(ctx context.Context, id int) ([]ProdCountResponse, error)
}

// service is a struct that represents a section service
type service struct {
	r Repository
}

// NewService returns a new instance of section service
func NewService(r Repository) Service {
	return &service{r}
}

// GetAll returns all sections
func (s *service) GetAll(ctx context.Context) (l []domain.Section, err error) {
	l, err = s.r.GetAll(ctx)
	return
}

// Get returns a sections given its id
func (s *service) Get(ctx context.Context, id int) (i domain.Section, err error) {
	i, err = s.r.Get(ctx, id)
	return
}

// Save stores a new section in the database
func (s *service) Save(ctx context.Context, sect domain.Section) (id int, err error) {
	// check if section already exists
	exists := s.r.Exists(ctx, sect.SectionNumber)
	if exists {
		err = ErrDuplicateSectNumber
		return
	}
	//save
	id, err = s.r.Save(ctx, sect)
	return
}

// Delete removes a section by its id from the database
func (s *service) Delete(ctx context.Context, id int) (err error) {
	err = s.r.Delete(ctx, id)
	return
}

// Update modifies fields of an existing section
func (s *service) Update(ctx context.Context, sect domain.Section) (err error) {
	// Check if section number already exists
	// - get unchanged section
	org, err := s.r.Get(ctx, sect.ID)
	if err != nil {
		return
	}
	// - if Section_Number is being updated, check if the new one already exists
	if sect.SectionNumber != org.SectionNumber {
		exists := s.r.Exists(ctx, sect.SectionNumber)
		if exists {
			err = ErrDuplicateSectNumber
			return
		}
	}
	// Save changes
	err = s.r.Update(ctx, sect)
	return
}

// Product count returns number of products for a given sections
// If no section is given, returns sesult for all sections
func (s *service) ProductCount(ctx context.Context, id int) (l []ProdCountResponse, err error) {
	l, err = s.r.ProductCount(ctx, id)
	return
}
