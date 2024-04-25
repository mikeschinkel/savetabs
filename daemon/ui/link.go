package ui

import (
	"html"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/safehtml"
)

type link struct {
	Id         int64
	rowId      int
	title      string
	url        *url.URL
	escapedURL string
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
	end := len(l.url.Host) - len(l.SecondLevelDomain())
	if end == 0 {
		goto end
	}
	sub = l.url.Host[0:end]
end:
	return sub
}

func (l link) SecondLevelDomain() (sld string) {
	segments := strings.Split(l.url.Host, ".")
	if len(segments) <= 2 {
		sld = l.url.Host
		goto end
	}
	sld = strings.Join(segments[len(segments)-2:], ".")
end:
	return sld
}

func (l link) ARIALabel() string {
	return "External Link: " + l.escapedURL
}
