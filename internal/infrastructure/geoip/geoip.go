package geoip

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
	"github.com/vondr/identity-go/internal/core"
)

type GeoIPService struct {
	db *geoip2.Reader
}

var geoIPService *GeoIPService

func InitGeoIP(dbPath string) error {
	if dbPath == "" {
		geoIPService = &GeoIPService{db: nil}
		return nil
	}

	db, err := geoip2.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open geoip database: %w", err)
	}

	geoIPService = &GeoIPService{db: db}
	return nil
}

func GetService() *GeoIPService {
	return geoIPService
}

func (s *GeoIPService) LookupCountry(ipStr string) (string, error) {
	if s.db == nil {
		return "", core.ErrGeoIPDisabled
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", core.ErrUnableToResolve
	}

	country, err := s.db.Country(ip)
	if err != nil {
		return "", fmt.Errorf("failed to lookup country: %w", err)
	}

	return country.Country.IsoCode, nil
}

func (s *GeoIPService) IsEnabled() bool {
	return s.db != nil
}

func (s *GeoIPService) IsPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.IsPrivate() || ip.IsLoopback()
}

func Close() {
	if geoIPService != nil && geoIPService.db != nil {
		geoIPService.db.Close()
	}
}
