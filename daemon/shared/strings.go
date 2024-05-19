package shared

import (
	"regexp"
	"strings"
)

var MatchAngleDelimited = regexp.MustCompile(`^<No (.+)>$`)

func Slugify(s string) (slug string) {

	if MatchAngleDelimited.MatchString(s) {
		// Replace `<No Whatever>` with `none` if angle brackets, otherwise just return input
		slug = "none"
		goto end
	}
	// Lowercase it
	slug = strings.ToLower(s)
	// Change all non-slug characters ([A-Zaz0-9-]) and convert to a dash (-)
	slug = MatchNonSlugCharacters.ReplaceAllString(slug, "-")
	// Condense all repeated dashes into a single dash
	slug = MatchDashes.ReplaceAllString(slug, "-")
end:
	return slug
}

type String struct {
	s string
}

func (s String) String() string {
	return s.s
}
