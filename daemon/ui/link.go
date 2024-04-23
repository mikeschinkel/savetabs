package ui

import (
	"bytes"
	"html"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/google/safehtml"
	"savetabs/sqlc"
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

var linksTemplate = GetTemplate("links")

func (v *Views) GetLinksHTML(ctx Context, host string, params FilterValueGetter, rawQuery string) (html []byte, status int, err error) {
	var out bytes.Buffer
	var ll []sqlc.ListFilteredLinksRow
	var links []link
	var ids []int64
	var linkIds []int64
	var values []string

	labels := []string{}
	for _, gt := range FilterTypes {
		values = params.GetFilterValues(gt)
		if len(values) == 0 {
			continue
		}
		// TODO: Change to using names for labels, not slugs
		labels = append(labels, params.GetFilterLabel(gt, strings.Join(values, ",")))
		switch gt {
		case MetaFilter:
			ids, err = v.Queries.ListLinkIdsByMetadata(ctx, sqlc.ListLinkIdsByMetadataParams{
				KvPairs:      values,
				LinkArchived: sqlc.NotArchived,
			})
		case GroupTypeFilter:
			ids, err = v.Queries.ListLinkIdsByGroupType(ctx, sqlc.ListLinkIdsByGroupTypeParams{
				GroupTypes:   values,
				LinkArchived: sqlc.NotArchived,
			})
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
		ll, err = v.Queries.ListFilteredLinks(ctx, sqlc.ListFilteredLinksParams{
			Ids:          linkIds,
			LinkArchived: sqlc.NotArchived,
		})
	}
	if err != nil {
		goto end
	}
	links = linksFromResources(ll)
	err = linksTemplate.Execute(&out, linkSet{
		apiURL:   makeURL(host),
		rawQuery: rawQuery,
		Links:    links,
		Label:    strings.Join(labels, " && "),
	})
	if err != nil {
		goto end
	}
	html = out.Bytes()
end:
	return html, http.StatusInternalServerError, err
}

func linksFromResources(ll []sqlc.ListFilteredLinksRow) (links []link) {
	links = make([]link, len(ll))
	for i, l := range ll {
		title := l.Title.String
		u, err := url.Parse(l.OriginalUrl)
		if err != nil {
			title = "ERROR: " + err.Error()
		}
		link := newLink(u)
		link.Id = l.ID
		link.rowId = i + 1
		link.title = title
		links[i] = link
	}
	return links
}
