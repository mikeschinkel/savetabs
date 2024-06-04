package ui

import (
	"html"
	"strconv"
	"strings"

	"github.com/google/safehtml"
	"savetabs/model"
	"savetabs/shared"
)

type htmlLinkArgs struct {
	Link     model.Link
	RowId    int
	DragDrop *dragDrop
}
type htmlLink struct {
	model.Link
	RowId int
	*dragDrop
}

func newHTMLLink(args htmlLinkArgs) htmlLink {
	return htmlLink{
		Link:     args.Link,
		RowId:    args.RowId,
		dragDrop: args.DragDrop,
	}
}

func (ll htmlLink) EscapedURL() string {
	return html.EscapeString(ll.Link.URL)
}

func (ll htmlLink) HTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`htmlLink`,
		strconv.FormatInt(ll.Link.Id, 10),
	)
}

func (ll htmlLink) DragDropId() safehtml.Identifier {
	return shared.MakeSafeIdf("%s-%d", ll.dragDrop, ll.Link.Id)
}

func (ll htmlLink) RowHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`links-row`,
		strconv.Itoa(ll.RowId),
	)
}

func (ll htmlLink) Domain() string {
	return ll.Scheme() +
		ll.Subdomain() +
		ll.SecondLevelDomain() +
		ll.TopLevelDomain() +
		ll.Port()
}

func (ll htmlLink) Title() (title string) {
	return ll.Link.Title
}

func (ll htmlLink) LocalTitle() (title string) {
	if strings.Contains(ll.Link.Title, "|") {
		title = strings.TrimSpace(ll.Link.Title[:strings.IndexByte(ll.Link.Title, '|')])
		goto end
	}
	title = ll.Link.Title
end:
	return title
}

func (ll htmlLink) Subdomain() (sub string) {
	if ll.Link.Subdomain == "" {
		goto end
	}
	sub = ll.Link.Subdomain + "."
end:
	return sub
}

func (ll htmlLink) Scheme() (s string) {
	if ll.Link.Scheme == "" {
		goto end
	}
	s = ll.Link.Scheme + "://"
end:
	return s
}

func (ll htmlLink) Path() string {
	return ll.Link.Path
}

func (ll htmlLink) Port() (port string) {
	if ll.Link.Port == "" {
		goto end
	}
	port = ":" + ll.Link.Port
end:
	return port
}

func (ll htmlLink) Query() (q string) {
	if ll.Link.Query == "" {
		goto end
	}
	q = "?" + ll.Link.Query
	if len(q) > 50 {
		q = "?<expand-to-see>"
	}
end:
	return q
}

func (ll htmlLink) Fragment() (frag string) {
	if ll.Link.Fragment == "" {
		goto end
	}
	frag = "#" + ll.Link.Fragment
	if len(frag) > 50 {
		frag = "#<expand-to-see>"
	}
end:
	return frag
}

func (ll htmlLink) TopLevelDomain() (tld string) {
	if ll.Link.TLD == "" {
		goto end
	}
	tld = "." + ll.Link.TLD
end:
	return tld
}

func (ll htmlLink) SecondLevelDomain() (sld string) {
	return ll.Link.SLD
}

func (ll htmlLink) ARIALabel() string {
	return "External Link: " + ll.EscapedURL()
}
