package shared

import (
	"errors"
)

var (
	ErrGroupTypeNotFoundForSlug   = errors.New("GroupType not found for slug")
	ErrGroupTypeNotFoundForLetter = errors.New("GroupType not found for letter")
	ErrGroupTypeIndexNotFound     = errors.New("GroupType index not found")
	ErrMenuTypeNotFound           = errors.New("MenuType not found")
	ErrContextMenuTypeNotFound    = errors.New("ContextMenuType not found")
	ErrMenuTypeIsNil              = errors.New("MenuType is nil")
	ErrInvalidFilterType          = errors.New("Invalid filter type")
	ErrInvalidGroupFilterFormat   = errors.New("Invalid group filter format")
	ErrInvalidMetaFilterFormat    = errors.New("Invalid meta filter format")
)
