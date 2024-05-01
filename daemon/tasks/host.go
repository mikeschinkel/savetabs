package tasks

import (
	"net/url"
	"regexp"
	"strings"
)

type Host struct {
	url          *url.URL
	Subdomain    string
	SLD          string
	Port         string
	IsIP         bool
	IsLocal      bool
	HasSubdomain bool
}

// TLD returns the Top Level Domain for the host e.g. `com`, `net`, `org`, `edu`
// TODO: Handle `co.uk`, etc.
func (h Host) TLD() (tld string) {
	var idx int
	var hn string

	if h.IsIP {
		// No TLD for an IP address
		goto end
	}
	if h.IsLocal {
		// No TLD for a local name like 'localhost' or 'my_mac'
		goto end
	}
	hn = h.url.Hostname()
	idx = strings.LastIndex(h.url.Hostname(), ".")
	tld = hn[idx+1:]
end:
	return tld
}

var matchIPv4Address = regexp.MustCompile(`^(((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4})$`)

// parseHost parses the Host property of a url.URL into convenient parts
// TODO: Handle `co.uk`, etc.
func parseHost(u *url.URL) (host Host) {
	host = Host{
		url:  u,
		Port: u.Port(),
	}
	h := u.Hostname()
	cnt := strings.Count(h, ".")
	switch {
	case cnt == 0:
		// When host is like 'localhost', or 'my_app'
		host.IsLocal = true
		host.SLD = h
	case cnt == 1:
		// No subdomain
		host.SLD, _, _ = strings.Cut(h, ".")
	case cnt == 3 && matchIPv4Address.MatchString(h):
		// Is an IP address
		host.IsIP = true
		host.SLD = h
	default:
		// Has subdomain(s)
		segments := strings.Split(h, ".")
		host.Subdomain = segments[len(segments)-3]
		host.SLD = segments[len(segments)-2 : len(segments)-1][0]
		host.HasSubdomain = true
	}
	return host
}
