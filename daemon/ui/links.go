package ui

import (
	"bytes"
	"fmt"
	"html"
	"strconv"

	"github.com/google/safehtml"
	"savetabs/sqlc"
)

type link struct {
	Id  int64
	URL string
}

func (l link) Identifier() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`link`,
		strconv.FormatInt(l.Id, 10),
	)
}
func (l link) Title() string {
	return html.EscapeString(l.URL)
}

func (l link) ARIALabel() string {
	return "External Link: " + html.EscapeString(l.URL)
}

type linkSet struct {
	apiURL string
	Links  []link
}

func (ls linkSet) HTMLLinksURL() string {
	return fmt.Sprintf("%s/html/links", ls.apiURL)
}

var linksTemplate = GetTemplate("links")

func GetLinksHTML(ctx Context, host string, params SlugsForGetter) (html []byte, err error) {
	var out bytes.Buffer
	var ll []sqlc.Link
	var links []link
	var ids []int64
	var linkIds []int64
	var slugs []string

	for _, gt := range GroupTypes {
		ch := string(gt)
		slugs = params.GetSlugsFor(ch)
		if len(slugs) == 0 {
			continue
		}
		if ch == MetaType {
			ids, err = queries.ListLinkIdsByMetadata(ctx, slugs)
		} else {
			ids, err = queries.ListLinkIdsByGroupSlugs(ctx, slugs)
		}
		if err != nil {
			goto end
		}
		if len(ids) == 0 {
			continue
		}
		linkIds = append(linkIds, ids...)
	}
	if len(linkIds) == 0 {
		ll, err = queries.ListLinks(ctx)
	} else {
		ll, err = queries.ListFilteredLinks(ctx, linkIds)
	}
	if err != nil {
		goto end
	}
	links = linksFromResources(ll)
	err = linksTemplate.Execute(&out, linkSet{
		apiURL: makeURL(host),
		Links:  links,
	})
	if err != nil {
		goto end
	}
	html = out.Bytes()
end:
	return html, err
}

func linksFromResources(ll []sqlc.Link) (links []link) {
	links = make([]link, len(ll))
	for i, l := range ll {
		links[i] = link{
			Id:  l.ID,
			URL: l.Url.String,
		}
	}
	return links
}
