package shared

import (
	"strings"

	"savetabs/regexes"
)

func Slugify(s string) string {
	slug := strings.ToLower(s)
	slug = regexes.MatchWhitespace.ReplaceAllString(slug, "-")
	slug = regexes.MatchNonSlugCharacters.ReplaceAllString(slug, "")
	return slug
}
