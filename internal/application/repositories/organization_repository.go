package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/vondr/identity-go/internal/core"
	"github.com/vondr/identity-go/internal/infrastructure/database/models"
)

type OrganizationRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	GetByHostname(ctx context.Context, hostname string) (*models.Organization, error)
	List(ctx context.Context) ([]*models.Organization, error)
	Create(ctx context.Context, org *models.Organization) error
	Update(ctx context.Context, org *models.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type GormOrganizationRepository struct {
	db *gorm.DB
}

func NewGormOrganizationRepository(db *gorm.DB) *GormOrganizationRepository {
	return &GormOrganizationRepository{db: db}
}

func (r *GormOrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	err := r.db.WithContext(ctx).First(&org, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return &org, nil
}

func (r *GormOrganizationRepository) GetByHostname(ctx context.Context, hostname string) (*models.Organization, error) {
	var org models.Organization
	err := r.db.WithContext(ctx).Where("hostname = ?", hostname).First(&org).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return &org, nil
}

func (r *GormOrganizationRepository) List(ctx context.Context) ([]*models.Organization, error) {
	var orgs []*models.Organization
	err := r.db.WithContext(ctx).Find(&orgs).Error
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *GormOrganizationRepository) Create(ctx context.Context, org *models.Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

func (r *GormOrganizationRepository) Update(ctx context.Context, org *models.Organization) error {
	return r.db.WithContext(ctx).Save(org).Error
}

func (r *GormOrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Organization{}, "id = ?", id).Error
}
