package shared

import (
	"regexp"
)

// MatchDashes matches all dashes ('-')
var MatchDashes = regexp.MustCompile(`-+`)

// MatchNonSlugCharacters matches only non-slug[*] characters.
// [*] Slug characters are letters (A-Z, a-z), numbers (0-9), and dashes ('-').
// The pattern [^A-Za-z0-9-] means:
// - ^ inside the brackets negates the character set, so it matches anything not in the set.
// - A-Za-z includes all uppercase and lowercase letters.
// - 0-9 includes all digits.
// - - is a literal dash, included in the character set.
var MatchNonSlugCharacters = regexp.MustCompile("[^A-Za-z0-9-]+")
