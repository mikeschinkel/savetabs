package sqlc

import (
	"strings"
)

func GroupTypeFromName(n string) (t string) {
	switch strings.ToLower(n) {
	case "category":
		t = "C"
	case "keyword":
		t = "K"
	case "tag":
		t = "T"
	case "tabgroup", "tab-group":
		t = "G"
	default:
		t = "I"
	}
	return t
}
