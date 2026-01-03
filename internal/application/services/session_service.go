package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/vondr/identity-go/internal/infrastructure/database"
	"github.com/vondr/identity-go/internal/infrastructure/database/models"
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService() *SessionService {
	return &SessionService{
		db: database.GetDB(),
	}
}

func (s *SessionService) RecordLogin(ctx context.Context, memberID uuid.UUID, email string, organizationID uuid.UUID, microsoftID string) error {
	session := &models.Session{
		MemberID:       memberID,
		Email:          email,
		OrganizationID: organizationID,
		MicrosoftID:    microsoftID,
	}

	return s.db.Create(session).Error
}

func (s *SessionService) GetLastLoginForMember(ctx context.Context, memberID uuid.UUID) (*time.Time, error) {
	var session models.Session
	err := s.db.Where("member_id = ?", memberID).Order("created_at DESC").First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session.CreatedAt, nil
}

func (s *SessionService) GetLastLoginsBatch(ctx context.Context, memberIDs []uuid.UUID) (map[uuid.UUID]*time.Time, error) {
	var sessions []models.Session
	err := s.db.Where("member_id IN ?", memberIDs).
		Select("DISTINCT ON (member_id) *").
		Order("member_id, created_at DESC").
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID]*time.Time)
	for _, session := range sessions {
		result[session.MemberID] = &session.CreatedAt
	}

	return result, nil
}
