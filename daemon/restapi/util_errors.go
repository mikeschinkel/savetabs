package restapi

import (
	"errors"
)

var (
	ErrValidatingHTTPRequest = errors.New("failed to validate HTTP request")
	ErrUnmarshallingJSON     = errors.New("failed to unmarshal JSON")
	ErrFailedUpsertLinks     = errors.New("failed to upsert links")
	ErrNoLinkIdsSubmitted    = errors.New("no link IDs submitted")
	ErrReadingHTTPBody       = errors.New("failed to read HTTP body")
)
