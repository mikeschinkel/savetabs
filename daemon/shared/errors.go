package shared

import (
	"fmt"
)

var (
	ErrGroupTypeNotFoundForSlug = fmt.Errorf("GroupType not found for slug")
	ErrGroupTypeNotFoundForType = fmt.Errorf("GroupType not found for type")
	ErrMenuTypeNotFound         = fmt.Errorf("MenuType not found")
	ErrMenuTypeIsNil            = fmt.Errorf("MenuType is nil")
)
