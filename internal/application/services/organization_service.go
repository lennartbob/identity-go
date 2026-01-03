package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/vondr/identity-go/internal/application/repositories"
	"github.com/vondr/identity-go/internal/infrastructure/database/models"
)

type OrganizationService struct {
	orgRepo repositories.OrganizationRepository
}

func NewOrganizationService(orgRepo repositories.OrganizationRepository) *OrganizationService {
	return &OrganizationService{
		orgRepo: orgRepo,
	}
}

func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	return s.orgRepo.GetByID(ctx, id)
}

func (s *OrganizationService) GetByHostname(ctx context.Context, hostname string) (*models.Organization, error) {
	return s.orgRepo.GetByHostname(ctx, hostname)
}

func (s *OrganizationService) List(ctx context.Context) ([]*models.Organization, error) {
	return s.orgRepo.List(ctx)
}

func (s *OrganizationService) Create(ctx context.Context, name string, hostname *string) (*models.Organization, error) {
	org := &models.Organization{
		Name:     name,
		Hostname: hostname,
	}

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *OrganizationService) Update(ctx context.Context, id uuid.UUID, name string, hostname *string) (*models.Organization, error) {
	org, err := s.orgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	org.Name = name
	if hostname != nil {
		org.Hostname = hostname
	}

	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *OrganizationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.orgRepo.Delete(ctx, id)
}
