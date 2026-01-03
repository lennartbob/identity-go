package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/vondr/identity-go/internal/core"
	"github.com/vondr/identity-go/internal/infrastructure/database/models"
)

type AppRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.App, error)
	GetByToken(ctx context.Context, token string) (*models.App, error)
	ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*models.App, error)
	Create(ctx context.Context, app *models.App) error
	Update(ctx context.Context, app *models.App) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type GormAppRepository struct {
	db *gorm.DB
}

func NewGormAppRepository(db *gorm.DB) *GormAppRepository {
	return &GormAppRepository{db: db}
}

func (r *GormAppRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.App, error) {
	var app models.App
	err := r.db.WithContext(ctx).First(&app, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return &app, nil
}

func (r *GormAppRepository) GetByToken(ctx context.Context, token string) (*models.App, error) {
	var app models.App
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&app).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return &app, nil
}

func (r *GormAppRepository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*models.App, error) {
	var apps []*models.App
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationID).Find(&apps).Error
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (r *GormAppRepository) Create(ctx context.Context, app *models.App) error {
	return r.db.WithContext(ctx).Create(app).Error
}

func (r *GormAppRepository) Update(ctx context.Context, app *models.App) error {
	return r.db.WithContext(ctx).Save(app).Error
}

func (r *GormAppRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.App{}, "id = ?", id).Error
}
