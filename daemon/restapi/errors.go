package restapi

import (
	"errors"

	"github.com/google/safehtml/template"
)

var (
	ErrFailedToUnmarshal = errors.New("failed to unmarshal JSON")
	ErrFailedUpsertLinks = errors.New("failed to upsert links")
)

var errorTemplate *template.Template

func SetErrorTemplate(t *template.Template) {
	errorTemplate = t
}
