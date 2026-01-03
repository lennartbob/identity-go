package services

import (
	"context"
	"errors"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/vondr/identity-go/internal/application/repositories"
	"github.com/vondr/identity-go/internal/infrastructure/database/models"
)

type AppAllowedCountryService interface {
	ListRules(ctx context.Context, appID uuid.UUID) ([]*models.AppAllowedCountry, error)
	ListCountryCodes(ctx context.Context, appID uuid.UUID) ([]string, error)
	AddCountry(ctx context.Context, appID uuid.UUID, countryCode string) (*models.AppAllowedCountry, error)
	RemoveCountry(ctx context.Context, appID uuid.UUID, countryCode string) error
	ReplaceCountries(ctx context.Context, appID uuid.UUID, countryCodes []string) ([]*models.AppAllowedCountry, error)
}

type AppAllowedCountryServiceImpl struct {
	repository repositories.AppAllowedCountryRepository
}

func NewAppAllowedCountryService(repository repositories.AppAllowedCountryRepository) *AppAllowedCountryServiceImpl {
	return &AppAllowedCountryServiceImpl{repository: repository}
}

func (s *AppAllowedCountryServiceImpl) ListRules(ctx context.Context, appID uuid.UUID) ([]*models.AppAllowedCountry, error) {
	return s.repository.ListByAppID(ctx, appID)
}

func (s *AppAllowedCountryServiceImpl) ListCountryCodes(ctx context.Context, appID uuid.UUID) ([]string, error) {
	rules, err := s.repository.ListByAppID(ctx, appID)
	if err != nil {
		return nil, err
	}
	codes := make([]string, len(rules))
	for i, rule := range rules {
		codes[i] = rule.CountryCode
	}
	return codes, nil
}

func (s *AppAllowedCountryServiceImpl) AddCountry(ctx context.Context, appID uuid.UUID, countryCode string) (*models.AppAllowedCountry, error) {
	normalized, err := s.normalizeCountryCode(countryCode)
	if err != nil {
		return nil, err
	}
	return s.repository.Add(ctx, appID, normalized)
}

func (s *AppAllowedCountryServiceImpl) RemoveCountry(ctx context.Context, appID uuid.UUID, countryCode string) error {
	normalized, err := s.normalizeCountryCode(countryCode)
	if err != nil {
		return err
	}
	return s.repository.Remove(ctx, appID, normalized)
}

func (s *AppAllowedCountryServiceImpl) ReplaceCountries(ctx context.Context, appID uuid.UUID, countryCodes []string) ([]*models.AppAllowedCountry, error) {
	normalized := make([]string, len(countryCodes))
	for i, code := range countryCodes {
		nc, err := s.normalizeCountryCode(code)
		if err != nil {
			return nil, err
		}
		normalized[i] = nc
	}
	return s.repository.Replace(ctx, appID, normalized)
}

func (s *AppAllowedCountryServiceImpl) normalizeCountryCode(value string) (string, error) {
	if value == "" {
		return "", errors.New("country code is required")
	}
	code := strings.TrimSpace(strings.ToUpper(value))
	if len(code) != 2 {
		return "", errors.New("country code must be a two-letter ISO code (e.g., 'NL')")
	}
	for _, r := range code {
		if !unicode.IsLetter(r) {
			return "", errors.New("country code must be a two-letter ISO code (e.g., 'NL')")
		}
	}
	return code, nil
}
