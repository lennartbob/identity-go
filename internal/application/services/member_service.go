package services

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/vondr/identity-go/internal/application/repositories"
	"github.com/vondr/identity-go/internal/core"
	"github.com/vondr/identity-go/internal/infrastructure/database/models"
)

type MemberService struct {
	memberRepo repositories.MemberRepository
	orgRepo    repositories.OrganizationRepository
}

func NewMemberService(memberRepo repositories.MemberRepository, orgRepo repositories.OrganizationRepository) *MemberService {
	return &MemberService{
		memberRepo: memberRepo,
		orgRepo:    orgRepo,
	}
}

func (s *MemberService) GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationMember, error) {
	return s.memberRepo.GetByID(ctx, id)
}

func (s *MemberService) GetByEmail(ctx context.Context, email string) (*models.OrganizationMember, error) {
	return s.memberRepo.GetByEmail(ctx, email)
}

func (s *MemberService) GetByMicrosoftID(ctx context.Context, microsoftID string) (*models.OrganizationMember, error) {
	return s.memberRepo.GetByMicrosoftID(ctx, microsoftID)
}

func (s *MemberService) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*models.OrganizationMember, error) {
	return s.memberRepo.ListByOrganizationID(ctx, organizationID)
}

func (s *MemberService) Create(ctx context.Context, member *models.OrganizationMember) error {
	return s.memberRepo.Create(ctx, member)
}

func (s *MemberService) Update(ctx context.Context, member *models.OrganizationMember) error {
	return s.memberRepo.Update(ctx, member)
}

func (s *MemberService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.memberRepo.Delete(ctx, id)
}

func (s *MemberService) CreateSystem(ctx context.Context, orgID uuid.UUID, orgName, email, microsoftID string, firstName, lastName *string) (*models.OrganizationMember, error) {
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil && err != core.ErrNotFound {
		return nil, err
	}

	if org == nil {
		org = &models.Organization{
			ID:       orgID,
			Name:     orgName,
			Hostname: nil,
		}
		if err := s.orgRepo.Create(ctx, org); err != nil {
			return nil, err
		}
	}

	member := &models.OrganizationMember{
		OrganizationID: orgID,
		MicrosoftID:    &microsoftID,
		Email:          strings.ToLower(email),
		FirstName:      firstName,
		LastName:       lastName,
		Role:           core.MemberRoleSystem,
	}

	if err := s.memberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *MemberService) LinkMicrosoftAccount(ctx context.Context, email, microsoftID string, firstName, lastName *string) (*models.OrganizationMember, error) {
	member, err := s.memberRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, core.ErrNotFound
	}

	member.MicrosoftID = &microsoftID
	if firstName != nil && member.FirstName == nil {
		member.FirstName = firstName
	}
	if lastName != nil && member.LastName == nil {
		member.LastName = lastName
	}

	if err := s.memberRepo.Update(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}
