package restapi

import (
	"errors"
)

var (
	ErrFailedToUnmarshal  = errors.New("failed to unmarshal JSON")
	ErrFailedUpsertLinks  = errors.New("failed to upsert links")
	ErrNoLinkIdsSubmitted = errors.New("no link IDs submitted")
)
