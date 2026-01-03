package models

import (
	"time"

	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
)

type StringArray []string

func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = StringArray{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, sa)
	case string:
		return json.Unmarshal([]byte(v), sa)
	default:
		return errors.New("unsupported type for StringArray")
	}
}

func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return "[]", nil
	}
	return json.Marshal(sa)
}

type App struct {
	ID              uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	OrganizationID  uuid.UUID   `gorm:"type:uuid;not null;index" json:"organization_id"`
	Name            string      `gorm:"type:varchar(255);not null" json:"name"`
	SubdomainLabels StringArray `gorm:"type:jsonb;default:'[]'" json:"subdomain_labels"`
	MainLabel       string      `gorm:"type:varchar(255);not null" json:"main_label"`
	Description     *string     `gorm:"type:text" json:"description"`
	CreatedAt       time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	IsPlatformApp   bool        `gorm:"type:boolean;not null;default:false" json:"is_platform_app"`
	Token           string      `gorm:"type:varchar(128);uniqueIndex;not null" json:"-"`
}
