package restapi

import (
	"errors"

	"github.com/google/safehtml/template"
)

var (
	ErrFailedToExtractKeyValues = errors.New("failed to extract key values")
	ErrFailedToUnmarshal        = errors.New("failed to unmarshal JSON")
	ErrFailedUpsertLinks        = errors.New("failed to upsert links")
	ErrFailedUpsertLinkGroups   = errors.New("failed to upsert link-groups")
	ErrFailedUpsertGroups       = errors.New("failed to upsert groups")
	ErrFailedUpsertKeyValues    = errors.New("failed to upsert key values")
	ErrUrlNotSpecified          = errors.New("url not specified")
	ErrLinkIsNil                = errors.New("link is nil")
	ErrUrlNotAbsolute           = errors.New("url is not absolute")
)

var errorTemplate *template.Template

func SetErrorTemplate(t *template.Template) {
	errorTemplate = t
}
