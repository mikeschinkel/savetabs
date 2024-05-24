package guard

import (
	"log/slog"
	"net/url"
	"strings"

	"savetabs/model"
	"savetabs/shared"
)

type AddLinksWithGroupsParams struct {
	Links []LinkWithGroup
}

type LinkWithGroup struct {
	URL       string
	Title     string
	GroupId   int64
	GroupType string
	Group     string
	parsedURL *url.URL
}

func (lwg *LinkWithGroup) Valid() (ok bool) {
	var err error

	lwg.parsedURL, err = url.Parse(lwg.URL)
	if err != nil {
		slog.Warn("Unable to parse URL", "link", lwg.URL, "error", err)
		goto end
	}
	ok = true
end:
	return ok
}

// Wanted returns true for Links/URLs we want to keep vs. ones we want to discard.
// TODO: Enhance to be end-user scriptable, ideally using two or more approaches:
//
//	https://github.com/expr-lang/expr | https://expr-lang.org/
//	https://github.com/google/cel-go
//	https://github.com/yuin/gopher-lua
//	https://github.com/dop251/goja
//	https://github.com/google/starlark-go
//	https://github.com/go-python/gpython
//	https://github.com/d5/tengo
//	https://github.com/mattn/anko
//	https://github.com/mikespook/goemphp
//	https://github.com/risor-io/risor
//	https://github.com/gentee/genteeextractDomains
//	https://code.google.com/archive/p/gotcl/
//	https://github.com/krotik/ecal
//	https://github.com/elsaland/elsa
//	https://github.com/antonvolkoff/goluajit
//	https://github.com/risor-io/risor
func (lwg *LinkWithGroup) Wanted() (ok bool) {

	switch lwg.URL {
	case "about:blank":
		goto end
	}

	switch lwg.parsedURL.Scheme { // TODO: Verify this has the values expected
	case "chrome", "edge", "view-source":
		goto end
	}

	ok = true
end:
	return ok
}

func (link *LinkWithGroup) EmptyURL() bool {
	return strings.TrimSpace(link.URL) == ""
}

func (link *LinkWithGroup) Sanitize() {
	if link.Group == "" {
		link.Group = "none"
	}
	if link.GroupType == "" {
		link.GroupType = shared.GroupTypeInvalid.Upper()
	}
}

func AddLinksWithGroupsIfNotExists(ctx Context, p AddLinksWithGroupsParams) (err error) {
	var gt shared.GroupType

	links := make([]model.LinkWithGroup, len(p.Links))
	for i, link := range p.Links {
		if link.EmptyURL() {
			slog.Warn("URL is unexpectedly empty", "link", link)
			continue
		}
		link.parsedURL, err = url.Parse(link.URL)
		if err != nil {
			slog.Warn("Unable to parse URL in link", "link", link)
			continue
		}
		if !link.Wanted() {
			continue
		}
		if !link.Valid() {
			continue
		}
		link.Sanitize()
		// TODO: Verify that link.GroupType is a slug and not a type
		gt, err = shared.ParseGroupTypeBySlug(link.GroupType)
		if err != nil {
			goto end
		}
		links[i] = model.LinkWithGroup{
			URL:         link.parsedURL,
			OriginalURL: link.URL,
			Title:       link.Title,
			GroupId:     link.GroupId,
			GroupType:   gt,
			Group:       link.Group,
		}
	}
	err = model.AddLinksWithGroupsIfNotExists(ctx, model.AddLinksWithGroupsParams{
		Links: links,
	})
end:
	return err
}
