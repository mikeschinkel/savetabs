package ui

import (
	"html"
	"strconv"
	"strings"

	"github.com/google/safehtml"
	"github.com/mikeschinkel/savetabs/daemon/model"
	"github.com/mikeschinkel/savetabs/daemon/shared"
)

type htmlLinkArgs struct {
	Link     model.Link
	RowId    int
	DragItem shared.DragDropItem
}

type htmlLink struct {
	model.Link
	RowId    int
	dragItem shared.DragDropItem
}

func newHTMLLink(args htmlLinkArgs) htmlLink {
	return htmlLink{
		Link:     args.Link,
		RowId:    args.RowId,
		dragItem: args.DragItem,
	}
}

func (hl htmlLink) EscapedURL() string {
	return html.EscapeString(hl.Link.URL)
}

func (hl htmlLink) HTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`htmlLink`,
		strconv.FormatInt(hl.Link.Id, 10),
	)
}

func (hl htmlLink) DragSources() safehtml.Identifier {
	return shared.MakeSafeId(hl.dragItem.DragSources())
}

func (hl htmlLink) DragParent() safehtml.Identifier {
	return shared.MakeSafeId(hl.dragItem.DragParent())
}

func (hl htmlLink) RowHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`links-row`,
		strconv.Itoa(hl.RowId),
	)
}

func (hl htmlLink) Domain() string {
	return hl.Scheme() +
		hl.Subdomain() +
		hl.SecondLevelDomain() +
		hl.TopLevelDomain() +
		hl.Port()
}

func (hl htmlLink) Title() (title string) {
	return hl.Link.Title
}

func (hl htmlLink) LocalTitle() (title string) {
	if strings.Contains(hl.Link.Title, "|") {
		title = strings.TrimSpace(hl.Link.Title[:strings.IndexByte(hl.Link.Title, '|')])
		goto end
	}
	title = hl.Link.Title
end:
	return title
}

func (hl htmlLink) Subdomain() (sub string) {
	if hl.Link.Subdomain == "" {
		goto end
	}
	sub = hl.Link.Subdomain + "."
end:
	return sub
}

func (hl htmlLink) Scheme() (s string) {
	if hl.Link.Scheme == "" {
		goto end
	}
	s = hl.Link.Scheme + "://"
end:
	return s
}

func (hl htmlLink) Path() string {
	return hl.Link.Path
}

func (hl htmlLink) Port() (port string) {
	if hl.Link.Port == "" {
		goto end
	}
	port = ":" + hl.Link.Port
end:
	return port
}

func (hl htmlLink) Query() (q string) {
	if hl.Link.Query == "" {
		goto end
	}
	q = "?" + hl.Link.Query
	if len(q) > 50 {
		q = "?<expand-to-see>"
	}
end:
	return q
}

func (hl htmlLink) Fragment() (frag string) {
	if hl.Link.Fragment == "" {
		goto end
	}
	frag = "#" + hl.Link.Fragment
	if len(frag) > 50 {
		frag = "#<expand-to-see>"
	}
end:
	return frag
}

func (hl htmlLink) TopLevelDomain() (tld string) {
	if hl.Link.TLD == "" {
		goto end
	}
	tld = "." + hl.Link.TLD
end:
	return tld
}

func (hl htmlLink) SecondLevelDomain() (sld string) {
	return hl.Link.SLD
}

func (hl htmlLink) ARIALabel() string {
	return "External Link: " + hl.EscapedURL()
}
