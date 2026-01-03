package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/vondr/identity-go/internal/core"
	"github.com/vondr/identity-go/internal/infrastructure/database/models"
)

type UserGroupRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.UserGroup, error)
	ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*models.UserGroup, error)
	ListByMemberID(ctx context.Context, memberID uuid.UUID) ([]*models.UserGroup, error)
	Create(ctx context.Context, group *models.UserGroup) error
	Update(ctx context.Context, group *models.UserGroup) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type GormUserGroupRepository struct {
	db *gorm.DB
}

func NewGormUserGroupRepository(db *gorm.DB) *GormUserGroupRepository {
	return &GormUserGroupRepository{db: db}
}

func (r *GormUserGroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.UserGroup, error) {
	var group models.UserGroup
	err := r.db.WithContext(ctx).First(&group, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return &group, nil
}

func (r *GormUserGroupRepository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*models.UserGroup, error) {
	var groups []*models.UserGroup
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationID).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *GormUserGroupRepository) ListByMemberID(ctx context.Context, memberID uuid.UUID) ([]*models.UserGroup, error) {
	var groups []*models.UserGroup
	err := r.db.WithContext(ctx).
		Joins("JOIN user_group_members ON user_groups.id = user_group_members.group_id").
		Where("user_group_members.member_id = ?", memberID).
		Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *GormUserGroupRepository) Create(ctx context.Context, group *models.UserGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *GormUserGroupRepository) Update(ctx context.Context, group *models.UserGroup) error {
	return r.db.WithContext(ctx).Save(group).Error
}

func (r *GormUserGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.UserGroup{}, "id = ?", id).Error
}
