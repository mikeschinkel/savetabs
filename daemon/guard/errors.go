package guard

import (
	"errors"
)

var (
	ErrInvalidUpsertAction = errors.New("invalid upsert action")
	ErrHTMLNotParsed       = errors.New("html not parsed")
	ErrUrlNotSpecified     = errors.New("url not specified")
	ErrInvalidKeyFormat    = errors.New("invalid key format (expected 'mi-<type>-<key>')")
)
