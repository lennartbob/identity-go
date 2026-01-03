package services

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/vondr/identity-go/internal/application/repositories"
	"github.com/vondr/identity-go/internal/core"
	"github.com/vondr/identity-go/internal/infrastructure/database/models"
)

type UserGroupService interface {
	List(ctx context.Context, organizationID uuid.UUID) ([]*models.UserGroup, error)
	Get(ctx context.Context, organizationID, groupID uuid.UUID) (*models.UserGroup, error)
	Create(ctx context.Context, organizationID uuid.UUID, name string, description *string, createdByMemberID *uuid.UUID) (*models.UserGroup, error)
	Update(ctx context.Context, organizationID, groupID uuid.UUID, name *string, description *string, updatedByMemberID *uuid.UUID) (*models.UserGroup, error)
	Delete(ctx context.Context, organizationID, groupID uuid.UUID) error
	ListMembers(ctx context.Context, organizationID, groupID uuid.UUID) ([]*models.UserGroupMember, error)
	AddMember(ctx context.Context, organizationID, groupID, memberID uuid.UUID) (*models.UserGroupMember, error)
	RemoveMember(ctx context.Context, organizationID, groupID, memberID uuid.UUID) error
	ListGroupsForMember(ctx context.Context, memberID uuid.UUID) ([]*models.UserGroup, error)
}

type UserGroupServiceImpl struct {
	groupRepo      repositories.UserGroupRepository
	memberRepo     repositories.UserGroupMemberRepository
	orgRepo        repositories.OrganizationRepository
	memberInfoRepo repositories.MemberRepository
}

func NewUserGroupService(
	groupRepo repositories.UserGroupRepository,
	memberRepo repositories.UserGroupMemberRepository,
	orgRepo repositories.OrganizationRepository,
	memberInfoRepo repositories.MemberRepository,
) *UserGroupServiceImpl {
	return &UserGroupServiceImpl{
		groupRepo:      groupRepo,
		memberRepo:     memberRepo,
		orgRepo:        orgRepo,
		memberInfoRepo: memberInfoRepo,
	}
}

func (s *UserGroupServiceImpl) normalizeName(name string) string {
	return strings.TrimSpace(name)
}

func (s *UserGroupServiceImpl) ensureOrgExists(ctx context.Context, organizationID uuid.UUID) error {
	if s.orgRepo == nil {
		return nil
	}
	_, err := s.orgRepo.GetByID(ctx, organizationID)
	if err != nil {
		return errors.New("organization not found")
	}
	return nil
}

func (s *UserGroupServiceImpl) ensureMemberBelongsToOrg(ctx context.Context, memberID, organizationID uuid.UUID) error {
	if s.memberInfoRepo == nil {
		return nil
	}
	member, err := s.memberInfoRepo.GetByID(ctx, memberID)
	if err != nil {
		return errors.New("member not found")
	}
	if member.DeletedAt.Valid {
		return errors.New("member not found")
	}
	if member.OrganizationID != organizationID {
		return errors.New("member belongs to a different organization")
	}
	return nil
}

func (s *UserGroupServiceImpl) List(ctx context.Context, organizationID uuid.UUID) ([]*models.UserGroup, error) {
	if err := s.ensureOrgExists(ctx, organizationID); err != nil {
		return nil, err
	}
	return s.groupRepo.ListByOrganizationID(ctx, organizationID)
}

func (s *UserGroupServiceImpl) Get(ctx context.Context, organizationID, groupID uuid.UUID) (*models.UserGroup, error) {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return nil, core.ErrNotFound
	}
	if group.OrganizationID != organizationID {
		return nil, core.ErrNotFound
	}
	return group, nil
}

func (s *UserGroupServiceImpl) Create(ctx context.Context, organizationID uuid.UUID, name string, description *string, createdByMemberID *uuid.UUID) (*models.UserGroup, error) {
	if err := s.ensureOrgExists(ctx, organizationID); err != nil {
		return nil, err
	}
	normalizedName := s.normalizeName(name)
	if normalizedName == "" {
		return nil, errors.New("group name cannot be empty")
	}
	if createdByMemberID != nil {
		if err := s.ensureMemberBelongsToOrg(ctx, *createdByMemberID, organizationID); err != nil {
			return nil, err
		}
	}
	group := &models.UserGroup{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		Name:           normalizedName,
		Description:    description,
	}
	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *UserGroupServiceImpl) Update(ctx context.Context, organizationID, groupID uuid.UUID, name *string, description *string, updatedByMemberID *uuid.UUID) (*models.UserGroup, error) {
	group, err := s.Get(ctx, organizationID, groupID)
	if err != nil {
		return nil, err
	}
	newName := group.Name
	if name != nil {
		newName = s.normalizeName(*name)
	}
	if newName == "" {
		return nil, errors.New("group name cannot be empty")
	}
	if updatedByMemberID != nil {
		if err := s.ensureMemberBelongsToOrg(ctx, *updatedByMemberID, organizationID); err != nil {
			return nil, err
		}
	}
	group.Name = newName
	if description != nil {
		group.Description = description
	}
	if err := s.groupRepo.Update(ctx, group); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *UserGroupServiceImpl) Delete(ctx context.Context, organizationID, groupID uuid.UUID) error {
	if _, err := s.Get(ctx, organizationID, groupID); err != nil {
		return err
	}
	return s.groupRepo.Delete(ctx, groupID)
}

func (s *UserGroupServiceImpl) ListMembers(ctx context.Context, organizationID, groupID uuid.UUID) ([]*models.UserGroupMember, error) {
	if _, err := s.Get(ctx, organizationID, groupID); err != nil {
		return nil, err
	}
	members, err := s.memberRepo.ListMembersByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	groupMembers := make([]*models.UserGroupMember, len(members))
	for i, member := range members {
		groupMembers[i] = &models.UserGroupMember{
			GroupID:  groupID,
			MemberID: member.ID,
			Member:   member,
		}
	}
	return groupMembers, nil
}

func (s *UserGroupServiceImpl) AddMember(ctx context.Context, organizationID, groupID, memberID uuid.UUID) (*models.UserGroupMember, error) {
	if _, err := s.Get(ctx, organizationID, groupID); err != nil {
		return nil, err
	}
	if err := s.ensureMemberBelongsToOrg(ctx, memberID, organizationID); err != nil {
		return nil, err
	}
	if err := s.memberRepo.AddMemberToGroup(ctx, groupID, memberID); err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"user_group_members_group_id_member_id_key\"" ||
			err.Error() == "pq: duplicate key value violates unique constraint \"user_group_members_group_id_member_id_key\"" {
			return nil, core.ErrConflict
		}
		return nil, err
	}
	member, err := s.memberInfoRepo.GetByID(ctx, memberID)
	if err != nil {
		return nil, err
	}
	return &models.UserGroupMember{
		GroupID:  groupID,
		MemberID: memberID,
		Member:   member,
	}, nil
}

func (s *UserGroupServiceImpl) RemoveMember(ctx context.Context, organizationID, groupID, memberID uuid.UUID) error {
	if _, err := s.Get(ctx, organizationID, groupID); err != nil {
		return err
	}
	return s.memberRepo.RemoveMemberFromGroup(ctx, groupID, memberID)
}

func (s *UserGroupServiceImpl) ListGroupsForMember(ctx context.Context, memberID uuid.UUID) ([]*models.UserGroup, error) {
	return s.groupRepo.ListByMemberID(ctx, memberID)
}
