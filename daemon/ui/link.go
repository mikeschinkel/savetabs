package ui

import (
	"html"
	"net/url"
	"strconv"

	"github.com/google/safehtml"
)

type link struct {
	Id         int64
	rowId      int
	title      string
	url        *url.URL
	escapedURL string
	scheme     string
	Host       string
	subdomain  string
	tld        string
	sld        string
	port       string
	path       string
	query      string
	fragment   string
}

func newLink(url *url.URL) link {
	l := link{
		url: url,
	}
	l.escapedURL = html.EscapeString(l.url.String())
	return l
}

func (l link) RowId() int {
	return l.rowId
}

func (l link) HTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`link`,
		strconv.FormatInt(l.Id, 10),
	)
}

func (l link) RowHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`links-row`,
		strconv.Itoa(l.rowId),
	)
}

func (l link) Domain() string {
	return l.Scheme() +
		l.Subdomain() +
		l.SecondLevelDomain() +
		l.TopLevelDomain() +
		l.Port()
}

func (l link) URL() string {
	return l.url.String()
}

func (l link) Title() (title string) {
	if l.title == "" {
		title = l.escapedURL
		goto end
	}
	title = l.title
end:
	return title
}

func (l link) Subdomain() (sub string) {
	if l.subdomain == "" {
		goto end
	}
	sub = l.subdomain + "."
end:
	return sub
}

func (l link) Scheme() (s string) {
	if l.scheme == "" {
		goto end
	}
	s = l.scheme + "://"
end:
	return s
}

func (l link) Path() string {
	return l.path
}

func (l link) Port() (port string) {
	if l.port == "" {
		goto end
	}
	port = ":" + l.port
end:
	return port
}

func (l link) Query() (q string) {
	if l.query == "" {
		goto end
	}
	q = "?" + l.query
	if len(q) > 50 {
		q = "?<expand-to-see>"
	}
end:
	return q
}

func (l link) Fragment() (frag string) {
	if l.fragment == "" {
		goto end
	}
	frag = "#" + l.fragment
	if len(frag) > 50 {
		frag = "#<expand-to-see>"
	}
end:
	return frag
}

func (l link) TopLevelDomain() (tld string) {
	if l.tld == "" {
		goto end
	}
	tld = "." + l.tld
end:
	return tld
}

func (l link) SecondLevelDomain() (sld string) {
	return l.sld
}

func (l link) ARIALabel() string {
	return "External Link: " + l.escapedURL
}
