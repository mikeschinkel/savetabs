package shared

import (
	"regexp"
	"strings"

	"savetabs/regexes"
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
	slug = regexes.MatchNonSlugCharacters.ReplaceAllString(slug, "-")
	// Condense all repeated dashes into a single dash
	slug = regexes.MatchDashes.ReplaceAllString(slug, "-")
end:
	return slug
}
