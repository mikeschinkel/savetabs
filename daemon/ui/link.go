package ui

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
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

func (v *Views) GetLinksHTML(ctx Context, host string, params FilterValueGetter) (html []byte, status int, err error) {
	var out bytes.Buffer
	var ll []sqlc.Link
	var links []link
	var ids []int64
	var linkIds []int64
	var values []string

	for _, gt := range FilterTypes {
		values = params.GetFilterValues(gt)
		if len(values) == 0 {
			continue
		}
		switch gt {
		case MetaFilter:
			ids, err = v.Queries.ListLinkIdsByMetadata(ctx, values)
		case GroupTypeFilter:
			ids, err = v.Queries.ListLinkIdsByGroupType(ctx, values)
		default:
			ids, err = v.Queries.ListLinkIdsByGroupSlugs(ctx, values)
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
		html = []byte("<div>No links for selection</div>")
		goto end
	} else {
		ll, err = v.Queries.ListFilteredLinks(ctx, linkIds)
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
	return html, http.StatusInternalServerError, err
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
