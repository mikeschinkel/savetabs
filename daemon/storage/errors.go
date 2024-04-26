package storage

import (
	"errors"
)

var (
	ErrDBNotANestedDCTX        = errors.New("database not a NestedDBTX")
	ErrFailedToUnmarshal       = errors.New("failed to unmarshal JSON")
	ErrFailedUpsertLinks       = errors.New("failed to upsert links")
	ErrFailedToExtractMetadata = errors.New("failed to extract metadata")
	ErrFailedUpsertLinkGroups  = errors.New("failed to upsert link-groups")
	ErrFailedUpsertGroups      = errors.New("failed to upsert groups")
	ErrFailedUpsertMetadata    = errors.New("failed to upsert metadata")
	ErrUrlNotSpecified         = errors.New("url not specified")
	ErrUrlNotAbsolute          = errors.New("url is not absolute")
	ErrFailedToArchiveLinks    = errors.New("failed to archive links")
)
