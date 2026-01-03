package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/vondr/identity-go/internal/infrastructure/database/models"
	"gorm.io/gorm"
)

type UserGroupMemberRepository interface {
	AddMemberToGroup(ctx context.Context, groupID, memberID uuid.UUID) error
	RemoveMemberFromGroup(ctx context.Context, groupID, memberID uuid.UUID) error
	ListMembersByGroupID(ctx context.Context, groupID uuid.UUID) ([]*models.OrganizationMember, error)
}

type GormUserGroupMemberRepository struct {
	db *gorm.DB
}

func NewGormUserGroupMemberRepository(db *gorm.DB) *GormUserGroupMemberRepository {
	return &GormUserGroupMemberRepository{db: db}
}

func (r *GormUserGroupMemberRepository) AddMemberToGroup(ctx context.Context, groupID, memberID uuid.UUID) error {
	member := &models.UserGroupMember{
		GroupID:  groupID,
		MemberID: memberID,
	}
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *GormUserGroupMemberRepository) RemoveMemberFromGroup(ctx context.Context, groupID, memberID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("group_id = ? AND member_id = ?", groupID, memberID).
		Delete(&models.UserGroupMember{}).Error
}

func (r *GormUserGroupMemberRepository) ListMembersByGroupID(ctx context.Context, groupID uuid.UUID) ([]*models.OrganizationMember, error) {
	var members []*models.OrganizationMember
	err := r.db.WithContext(ctx).
		Joins("JOIN user_group_members ON organization_members.id = user_group_members.member_id").
		Where("user_group_members.group_id = ?", groupID).
		Find(&members).Error
	if err != nil {
		return nil, err
	}
	return members, nil
}
