package models

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OrganizationID uuid.UUID  `gorm:"type:uuid;not null;index" json:"organization_id"`
	Token          string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	Email          string     `gorm:"type:varchar(320);not null" json:"email"`
	Role           string     `gorm:"type:varchar(20);not null" json:"role"`
	ExpiresAt      time.Time  `gorm:"not null" json:"expires_at"`
	AcceptedAt     *time.Time `json:"accepted_at"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (i *Invitation) TableName() string {
	return "invitations"
}
