package guard

import (
	"errors"
)

var (
	ErrInvalidUpsertAction   = errors.New("invalid upsert action")
	ErrHTMLNotParsed         = errors.New("html not parsed")
	ErrUrlNotSpecified       = errors.New("url not specified")
	ErrInvalidMenuItemFormat = errors.New("invalid menu item format (expected '<type>-<key>')")
)
