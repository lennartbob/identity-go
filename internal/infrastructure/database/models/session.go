package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	MemberID       uuid.UUID `gorm:"type:uuid;not null;index" json:"member_id"`
	Email          string    `gorm:"type:varchar(320);not null" json:"email"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null" json:"organization_id"`
	MicrosoftID    string    `gorm:"type:varchar(255);not null" json:"microsoft_id"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}
