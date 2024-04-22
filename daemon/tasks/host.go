package tasks

import (
	"net/url"
	"regexp"
	"strings"
)

type Host struct {
	url          *url.URL
	Sld          string
	Port         string
	IsIP         bool
	IsLocal      bool
	HasSubdomain bool
}

func (h Host) Subdomain() (sub string) {
	var host string
	var end int

	if h.IsIP {
		goto end
	}
	host = h.url.Hostname()
	end = len(host) - len(h.Sld)
	if end == 0 {
		goto end
	}
	sub = host[0 : end-1]
end:
	return sub
}

// TLD returns the Top Level Domain for the host e.g. `com`, `net`, `org`, `edu`
// TODO: Handle `co.uk`, etc.
func (h Host) TLD() (tld string) {
	var idx int
	if h.IsIP {
		// No TLD for an IP address
		goto end
	}
	if h.IsLocal {
		// No TLD for a local name like 'localhost' or 'my_mac'
		goto end
	}
	idx = strings.LastIndex(h.Sld, ".")
	tld = h.Sld[idx+1:]
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
		host.Sld = h
	case cnt == 1:
		// No subdomain
		host.Sld, _, _ = strings.Cut(h, ".")
	case cnt == 3 && matchIPv4Address.MatchString(h):
		// Is an IP address
		host.IsIP = true
		host.Sld = h
	default:
		// Has subdomain(s)
		segments := strings.Split(h, ".")
		host.Sld = segments[len(segments)-2]
		host.HasSubdomain = true
	}
	return host
}
