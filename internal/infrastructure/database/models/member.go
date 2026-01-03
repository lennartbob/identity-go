package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/vondr/identity-go/internal/core"
)

type OrganizationMember struct {
	ID             uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OrganizationID uuid.UUID       `gorm:"type:uuid;not null;index" json:"organization_id"`
	Organization   *Organization   `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	MicrosoftID    *string         `gorm:"type:varchar(255);uniqueIndex" json:"microsoft_id"`
	Email          string          `gorm:"type:varchar(320);not null" json:"email"`
	FirstName      *string         `gorm:"type:varchar(200)" json:"first_name"`
	LastName       *string         `gorm:"type:varchar(200)" json:"last_name"`
	Role           core.MemberRole `gorm:"type:varchar(20);not null;default:'member'" json:"role"`
	CreatedAt      time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (m *OrganizationMember) TableName() string {
	return "organization_members"
}
