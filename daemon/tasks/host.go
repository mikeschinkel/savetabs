package tasks

import (
	"net/url"
)

type LinkToParse struct {
	url          *url.URL
	Subdomain    string
	SLD          string
	Port         string
	IsIP         bool
	IsLocal      bool
	HasSubdomain bool
}
