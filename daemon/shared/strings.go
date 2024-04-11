package shared

import (
	"strings"

	"savetabs/regexes"
)

func Slugify(s string) string {
	slug := strings.ToLower(s)
	// Change all non-slug characters ([A-Zaz0-9-]) and convert to a dash (-)
	slug = regexes.MatchNonSlugCharacters.ReplaceAllString(slug, "-")
	// Condense all repeated dashes into a single dash
	slug = regexes.MatchDashes.ReplaceAllString(slug, "-")
	if slug != "-none-" {
		// Trim dashes from either end, unless `s` was originally `<none>`
		slug = strings.Trim(slug, "-")
	}
	return slug
}
