package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vondr/identity-go/internal/application/repositories"
	"github.com/vondr/identity-go/internal/infrastructure/database/models"
)

type AppService struct {
	appRepo repositories.AppRepository
	orgRepo repositories.OrganizationRepository
}

func NewAppService(appRepo repositories.AppRepository, orgRepo repositories.OrganizationRepository) *AppService {
	return &AppService{
		appRepo: appRepo,
		orgRepo: orgRepo,
	}
}

func (s *AppService) GetByID(ctx context.Context, id uuid.UUID) (*models.App, error) {
	return s.appRepo.GetByID(ctx, id)
}

func (s *AppService) GetByToken(ctx context.Context, token string) (*models.App, error) {
	return s.appRepo.GetByToken(ctx, token)
}

func (s *AppService) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*models.App, error) {
	return s.appRepo.ListByOrganizationID(ctx, organizationID)
}

func (s *AppService) Create(ctx context.Context, app *models.App) error {
	return s.appRepo.Create(ctx, app)
}

func (s *AppService) Update(ctx context.Context, app *models.App) error {
	return s.appRepo.Update(ctx, app)
}

func (s *AppService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.appRepo.Delete(ctx, id)
}

func (s *AppService) GetAllowedDomainsForOrganization(ctx context.Context, organizationID uuid.UUID, hostname *string) ([]string, error) {
	apps, err := s.appRepo.ListByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	domains := make(map[string]bool)

	if hostname != nil && *hostname != "" {
		domains[*hostname] = true
	}

	for _, app := range apps {
		for _, label := range app.SubdomainLabels {
			if label == app.MainLabel {
				if hostname != nil && *hostname != "" {
					domain := fmt.Sprintf("%s.%s", label, *hostname)
					domains[domain] = true
				}
			} else {
				if hostname != nil && *hostname != "" {
					domain := fmt.Sprintf("%s.%s", label, *hostname)
					domains[domain] = true
				}
			}
		}
	}

	result := make([]string, 0, len(domains))
	for domain := range domains {
		result = append(result, domain)
	}

	return result, nil
}

func (s *AppService) GetDomainAppMap(ctx context.Context, organizationID uuid.UUID, hostname *string) (map[string]*models.App, error) {
	apps, err := s.appRepo.ListByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*models.App)

	for _, app := range apps {
		for _, label := range app.SubdomainLabels {
			var domain string
			if label == app.MainLabel && hostname != nil && *hostname != "" {
				domain = fmt.Sprintf("%s.%s", label, *hostname)
			} else if hostname != nil && *hostname != "" {
				domain = fmt.Sprintf("%s.%s", label, *hostname)
			} else {
				domain = label
			}
			result[domain] = app
		}
	}

	return result, nil
}

func (s *AppService) GenerateToken() string {
	return uuid.New().String()
}
