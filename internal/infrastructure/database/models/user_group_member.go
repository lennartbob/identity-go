package models

import (
	"github.com/google/uuid"
)

type UserGroupMember struct {
	ID       uuid.UUID           `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	GroupID  uuid.UUID           `gorm:"type:uuid;not null;index" json:"group_id"`
	MemberID uuid.UUID           `gorm:"type:uuid;not null;index" json:"member_id"`
	Group    *UserGroup          `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Member   *OrganizationMember `gorm:"foreignKey:MemberID" json:"member,omitempty"`
}

func (ugm *UserGroupMember) TableName() string {
	return "user_group_members"
}
