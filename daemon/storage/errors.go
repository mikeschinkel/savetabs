package storage

import (
	"errors"
)

var (
	ErrFailedToUnmarshal      = errors.New("failed to unmarshal JSON")
	ErrFailedUpsertLinks      = errors.New("failed to upsert links")
	ErrFailedUpsertLink       = errors.New("failed to upsert link")
	ErrFailedUpsertLinkMeta   = errors.New("failed to upsert link meta")
	ErrFailedLoadLinkByUrl    = errors.New("failed to load link by URL")
	ErrFailedToExtractMeta    = errors.New("failed to extract meta")
	ErrFailedUpsertLinkGroups = errors.New("failed to upsert link-groups")
	ErrFailedUpsertGroups     = errors.New("failed to upsert groups")
	ErrFailedUpsertMeta       = errors.New("failed to upsert meta")
	ErrUrlNotSpecified        = errors.New("url not specified")
	ErrUrlNotAbsolute         = errors.New("url is not absolute")
	ErrFailedToArchiveLinks   = errors.New("failed to archive links")
	ErrFoundInLink            = errors.New("error found in link")
	ErrHTMLNotParsed          = errors.New("html not parsed")
)
