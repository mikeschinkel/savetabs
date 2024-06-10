package shared

import (
	"regexp"
)

// safeIdentifierPattern matches Identifiers that are a valid HTML5 'id' attribute.
//
//	This pattern allows:
//	- ASCII letters (a-z, A-Z)
//	- Digits (0-9)
//	- Common special characters: hyphens (-), underscores (_), colons (:), and periods (.)
//	- Unicode characters above U+007F to support internationalization.
var matchIdentifier = regexp.MustCompile(`^[\p{L}\p{N}_\-.:]+$`)

type Identifier struct {
	id    string
	Valid bool
}

func (id Identifier) String() string {
	return id.id
}

func NewIdentifier(id string) Identifier {
	identifier := Identifier{
		id: matchIdentifier.FindString(id),
	}
	if id != identifier.id {
		identifier.Valid = false
	}
	return identifier
}
