package shared

import (
	"errors"
)

var (
	ErrGroupTypeNotFoundForSlug   = errors.New("GroupType not found for slug")
	ErrGroupTypeNotFoundForLetter = errors.New("GroupType not found for letter")
	ErrMenuTypeNotFound           = errors.New("MenuType not found")
	ErrMenuTypeIsNil              = errors.New("MenuType is nil")
	ErrInvalidFilterType          = errors.New("Invalid filter type")
	ErrInvalidGroupFilterFormat   = errors.New("Invalid group filter format")
	ErrInvalidMetaFilterFormat    = errors.New("Invalid meta filter format")
)
