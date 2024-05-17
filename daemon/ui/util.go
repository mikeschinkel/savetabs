package ui

import (
	"fmt"
	"net/url"
)

func panicf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

func httpFormGet(form url.Values, key string) []string {
	items, ok := form[key]
	if !ok {
		items = []string{}
	}
	return items
}
