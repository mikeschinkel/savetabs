package shared

import (
	"fmt"
)

var (
	ErrGroupTypeNotFoundForSlug = fmt.Errorf("GroupType not found for slug")
	ErrGroupTypeNotFoundForCode = fmt.Errorf("GroupType not found for code")
)
