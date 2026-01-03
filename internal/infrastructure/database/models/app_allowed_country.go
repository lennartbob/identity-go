package models

import (
	"github.com/google/uuid"
)

type AppAllowedCountry struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AppID       uuid.UUID `gorm:"type:uuid;not null;index" json:"app_id"`
	App         *App      `gorm:"foreignKey:AppID" json:"app,omitempty"`
	CountryCode string    `gorm:"type:varchar(2);not null;index" json:"country_code"`
}

func (aac *AppAllowedCountry) TableName() string {
	return "app_allowed_countries"
}
