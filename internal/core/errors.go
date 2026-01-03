package core

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrConflict        = errors.New("resource already exists")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrBadRequest      = errors.New("bad request")
	ErrInvalidToken    = errors.New("invalid token")
	ErrExpiredToken    = errors.New("expired token")
	ErrInvalidSession  = errors.New("invalid session")
	ErrNoInvitation    = errors.New("no invitation found")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidDomain   = errors.New("invalid domain")
	ErrInvalidCountry  = errors.New("invalid country")
	ErrGeoIPDisabled   = errors.New("geoip not configured")
	ErrUnableToResolve = errors.New("unable to resolve")
)
