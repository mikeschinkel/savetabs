package ui

import (
	"fmt"
)

func makeURL(host string) string {
	return "http://" + host
}

func panicf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}
