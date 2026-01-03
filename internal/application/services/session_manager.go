package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vondr/identity-go/internal/infrastructure/cache"
)

type SessionManager struct {
	sessionRepo cache.SessionRepository
}

func NewSessionManager(sessionRepo cache.SessionRepository) *SessionManager {
	return &SessionManager{
		sessionRepo: sessionRepo,
	}
}

func (s *SessionManager) CreateSession(ctx context.Context, memberID uuid.UUID, email string, organizationID uuid.UUID, microsoftID string) (string, error) {
	token := uuid.New().String()
	ttl := time.Duration(7*24) * time.Hour

	sessionData := cache.SessionData{
		MemberID:       memberID.String(),
		Email:          email,
		OrganizationID: organizationID.String(),
		MicrosoftID:    microsoftID,
	}

	if err := s.sessionRepo.CreateSession(ctx, token, sessionData, ttl); err != nil {
		return "", err
	}

	return token, nil
}

func (s *SessionManager) GetSession(ctx context.Context, token string) (*cache.SessionData, error) {
	return s.sessionRepo.GetSession(ctx, token)
}

func (s *SessionManager) DeleteSession(ctx context.Context, token string) error {
	return s.sessionRepo.DeleteSession(ctx, token)
}
