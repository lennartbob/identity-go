package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/vondr/identity-go/internal/core"
	"github.com/vondr/identity-go/internal/infrastructure/database/models"
)

type MemberRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationMember, error)
	GetByEmail(ctx context.Context, email string) (*models.OrganizationMember, error)
	GetByMicrosoftID(ctx context.Context, microsoftID string) (*models.OrganizationMember, error)
	ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*models.OrganizationMember, error)
	Create(ctx context.Context, member *models.OrganizationMember) error
	Update(ctx context.Context, member *models.OrganizationMember) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type GormMemberRepository struct {
	db *gorm.DB
}

func NewGormMemberRepository(db *gorm.DB) *GormMemberRepository {
	return &GormMemberRepository{db: db}
}

func (r *GormMemberRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationMember, error) {
	var member models.OrganizationMember
	err := r.db.WithContext(ctx).First(&member, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return &member, nil
}

func (r *GormMemberRepository) GetByEmail(ctx context.Context, email string) (*models.OrganizationMember, error) {
	var member models.OrganizationMember
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return &member, nil
}

func (r *GormMemberRepository) GetByMicrosoftID(ctx context.Context, microsoftID string) (*models.OrganizationMember, error) {
	var member models.OrganizationMember
	err := r.db.WithContext(ctx).Where("microsoft_id = ?", microsoftID).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return &member, nil
}

func (r *GormMemberRepository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*models.OrganizationMember, error) {
	var members []*models.OrganizationMember
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationID).Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (r *GormMemberRepository) Create(ctx context.Context, member *models.OrganizationMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *GormMemberRepository) Update(ctx context.Context, member *models.OrganizationMember) error {
	return r.db.WithContext(ctx).Save(member).Error
}

func (r *GormMemberRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.OrganizationMember{}, "id = ?", id).Error
}
