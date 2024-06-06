package restapi

import (
	"errors"
)

var (
	ErrFailedToFixThisTODO = errors.New("failed to fix this TODO")
	ErrFailedToUnmarshal   = errors.New("failed to unmarshal JSON")
	ErrFailedUpsertLinks   = errors.New("failed to upsert links")
	ErrNoLinkIdsSubmitted  = errors.New("no link IDs submitted")
)
