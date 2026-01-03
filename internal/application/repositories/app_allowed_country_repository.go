package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/vondr/identity-go/internal/core"
	"github.com/vondr/identity-go/internal/infrastructure/database/models"
	"gorm.io/gorm"
)

type AppAllowedCountryRepository interface {
	ListByAppID(ctx context.Context, appID uuid.UUID) ([]*models.AppAllowedCountry, error)
	Add(ctx context.Context, appID uuid.UUID, countryCode string) (*models.AppAllowedCountry, error)
	Remove(ctx context.Context, appID uuid.UUID, countryCode string) error
	Replace(ctx context.Context, appID uuid.UUID, countryCodes []string) ([]*models.AppAllowedCountry, error)
}

type GormAppAllowedCountryRepository struct {
	db *gorm.DB
}

func NewGormAppAllowedCountryRepository(db *gorm.DB) *GormAppAllowedCountryRepository {
	return &GormAppAllowedCountryRepository{db: db}
}

func (r *GormAppAllowedCountryRepository) ListByAppID(ctx context.Context, appID uuid.UUID) ([]*models.AppAllowedCountry, error) {
	var countries []*models.AppAllowedCountry
	err := r.db.WithContext(ctx).
		Where("app_id = ?", appID).
		Order("country_code").
		Find(&countries).Error
	if err != nil {
		return nil, err
	}
	return countries, nil
}

func (r *GormAppAllowedCountryRepository) Add(ctx context.Context, appID uuid.UUID, countryCode string) (*models.AppAllowedCountry, error) {
	country := &models.AppAllowedCountry{
		ID:          uuid.New(),
		AppID:       appID,
		CountryCode: countryCode,
	}
	err := r.db.WithContext(ctx).Create(country).Error
	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"app_allowed_countries_app_id_country_code_key\"" ||
			err.Error() == "pq: duplicate key value violates unique constraint \"app_allowed_countries_app_id_country_code_key\"" {
			return nil, core.ErrConflict
		}
		return nil, err
	}
	return country, nil
}

func (r *GormAppAllowedCountryRepository) Remove(ctx context.Context, appID uuid.UUID, countryCode string) error {
	result := r.db.WithContext(ctx).
		Where("app_id = ? AND country_code = ?", appID, countryCode).
		Delete(&models.AppAllowedCountry{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return core.ErrNotFound
	}
	return nil
}

func (r *GormAppAllowedCountryRepository) Replace(ctx context.Context, appID uuid.UUID, countryCodes []string) ([]*models.AppAllowedCountry, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("app_id = ?", appID).Delete(&models.AppAllowedCountry{}).Error; err != nil {
			return err
		}

		if len(countryCodes) > 0 {
			var countries []*models.AppAllowedCountry
			for _, code := range countryCodes {
				countries = append(countries, &models.AppAllowedCountry{
					ID:          uuid.New(),
					AppID:       appID,
					CountryCode: code,
				})
			}
			if err := tx.Create(&countries).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.ListByAppID(ctx, appID)
}
