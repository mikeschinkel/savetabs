package restapi

import (
	"errors"

	"github.com/google/safehtml/template"
)

var (
	ErrFailedToExtractMetadata = errors.New("failed to extract metadata")
	ErrFailedToUnmarshal       = errors.New("failed to unmarshal JSON")
	ErrFailedUpsertLinks       = errors.New("failed to upsert links")
	ErrFailedUpsertLinkGroups  = errors.New("failed to upsert link-groups")
	ErrFailedUpsertGroups      = errors.New("failed to upsert groups")
	ErrFailedUpsertMetadata    = errors.New("failed to upsert metadata")
	ErrUrlNotSpecified         = errors.New("url not specified")
	ErrUrlNotAbsolute          = errors.New("url is not absolute")
)

var errorTemplate *template.Template

func SetErrorTemplate(t *template.Template) {
	errorTemplate = t
}
