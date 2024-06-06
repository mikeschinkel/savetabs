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
	ErrInvalidFilterType          = errors.New("invalid filter type")
	ErrInvalidGroupFilterFormat   = errors.New("invalid group filter format")
	ErrInvalidMetaFilterFormat    = errors.New("invalid meta filter format")
	ErrParamIsNotKindOfStruct     = errors.New("input is not a kind of struct")
	ErrNonZeroFieldValueRequired  = errors.New("non-zero field value required")
)
