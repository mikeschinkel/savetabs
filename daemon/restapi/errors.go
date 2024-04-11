package restapi

import (
	"errors"

	"github.com/google/safehtml/template"
)

var (
	ErrFailedToExtractKeyValues   = errors.New("failed to extract key values")
	ErrFailedToUnmarshal          = errors.New("failed to unmarshal JSON")
	ErrFailedUpsertResources      = errors.New("failed to upsert resources")
	ErrFailedUpsertResourceGroups = errors.New("failed to upsert resource-groups")
	ErrFailedUpsertGroups         = errors.New("failed to upsert groups")
	ErrFailedUpsertKeyValues      = errors.New("failed to upsert key values")
	ErrUrlNotSpecified            = errors.New("url not specified")
	ErrResourceIsNil              = errors.New("resource is nil")
	ErrUrlNotAbsolute             = errors.New("url is not absolute")
)

var errorTemplate *template.Template

func SetErrorTemplate(t *template.Template) {
	errorTemplate = t
}
