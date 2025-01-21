package model

import (
	"log/slog"
	"net/url"
	"regexp"
	"strings"

	"github.com/mikeschinkel/savetabs/daemon/shared"
	"github.com/mikeschinkel/savetabs/daemon/storage"
)

type UnparsedLink storage.UnparsedLink

type ParsedLink struct {
	url          *url.URL
	IsIP         bool
	IsLocal      bool
	HasSubdomain bool
	valid        bool
	storage.ParsedLink
}

func (l *ParsedLink) IsValid() bool {
	return l.valid

}

var matchIPv4Address = regexp.MustCompile(`^(((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4})$`)

// NewParsedLink returns a new ParsedLink given an UnparsedLink.
// TODO: Handle `co.uk`, etc.
func NewParsedLink(ul UnparsedLink) (link ParsedLink) {
	var u *url.URL
	var hn string
	var cnt int

	// Check pre-conditions for the parameter.
	// Given the app's architecture, these should never fail, but if they do it is a programming error.
	if ul.OriginalURL == "" {
		slog.Error("UnparsedLink.OriginalURL should not be empty.")
		goto end
	}
	u = ul.URL
	if u == nil {
		slog.Error("UnparsedLink.URL is nil, but should contain an instantiated *url.URL.",
			"url", ul.OriginalURL,
		)
		goto end
	}
	link.valid = true
	link = ParsedLink{
		url: u,
		ParsedLink: storage.ParsedLink{
			Title:       ul.Title,
			Scheme:      u.Scheme,
			OriginalURL: ul.OriginalURL,
			Port:        u.Port(),
			Path:        u.Path,
			Query:       strings.ReplaceAll(u.RawQuery, "%20", "+"),
			Fragment:    strings.ReplaceAll(u.Fragment, "%20", "+"),
		},
	}
	hn = u.Hostname()
	cnt = strings.Count(hn, ".")
	switch {
	case cnt == 0:
		// When link is like 'localhost', or 'my_app'
		link.IsLocal = true
		link.ParsedLink.SLD = hn
	case cnt == 1:
		// No subdomain
		link.ParsedLink.SLD, _, _ = strings.Cut(hn, ".")
	case cnt == 3 && matchIPv4Address.MatchString(hn):
		// Is an IP address
		link.IsIP = true
		link.ParsedLink.SLD = hn
	default:
		// Has subdomain(s)
		segments := strings.Split(hn, ".")
		link.ParsedLink.Subdomain = segments[len(segments)-3]
		link.ParsedLink.SLD = segments[len(segments)-2 : len(segments)-1][0]
		link.HasSubdomain = true
	}
	// Finally, derive TLD
	if link.IsIP {
		// No TLD for an IP address
		goto end
	}
	if link.IsLocal {
		// No TLD for a local name like 'localhost' or 'my_mac'
		goto end
	}
	link.ParsedLink.TLD = hn[strings.LastIndex(hn, ".")+1:]
end:
	return link
}

func LatestUnparsedLinks(ctx Context) (links []UnparsedLink, err error) {
	var storageLinks []storage.UnparsedLink

	storageLinks, err = storage.LatestUnparsedLinks(ctx)
	if err != nil {
		goto end
	}
	links = shared.ConvertSlice(storageLinks, func(link storage.UnparsedLink) UnparsedLink {
		return UnparsedLink(link)
	})
end:
	return links, err
}

func UpdateUnparsedLink(ctx Context, link UnparsedLink) error {
	pl := NewParsedLink(link)
	return storage.UpdateUnparsedLink(ctx, pl.ParsedLink)
}
