package ui

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/google/safehtml"
	"savetabs/sqlc"
)

type linkSet struct {
	apiURL     string
	Links      []link
	Label      string
	requestURI string
	queryJSON  string
}

func (ls linkSet) HeaderHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`links-row-head`)
}
func (ls linkSet) FooterHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`links-row-foot`)
}
func (ls linkSet) URLQuery() safehtml.URL {
	parts := strings.Split(ls.requestURI+"?", "?")
	return safehtml.URLSanitized("?" + parts[1])
}
func (ls linkSet) QueryJSON() safehtml.JSON {
	j, err := safehtml.JSONFromValue(ls.queryJSON)
	if err != nil {
		slog.Error("Unable to create safe JSON", "json", ls.queryJSON)
		j = safehtml.JSONFromConstant("{}")
	}
	return j
}
func (ls linkSet) NumLinks() int {
	return len(ls.Links)
}
func (ls linkSet) HTMLLinksURL() string {
	return fmt.Sprintf("%s/html/linkset", ls.apiURL)
}
func (ls linkSet) TableHeaderFooterHTML() safehtml.HTML {
	return safehtml.HTMLFromConstant(`
<th class="p-0.5">#</th>
<th class="p-0.5">
	<label>
		<input type="checkbox" @change="maybeConfirmCheckAll" class="check-all"> 
	</label>
</th>
<th class="p-0.5 text-center">Link</th>
<th class="p-0.5 text-right">Sub</th>
<th class="p-0.5">Domain</th>
<th class="p-0.5">Title</th>
`)
}

var linkSetTemplate = GetTemplate("link-set")

func (v *Views) GetLinkSetHTML(ctx Context, host, requestURI string, params FilterGetter) (html safehtml.HTML, status int, err error) {
	var ll []sqlc.ListFilteredLinksRow
	var links []link
	var ids []int64
	var linkIds []int64
	var values []string
	var queryJSON string

	for _, gt := range FilterTypes {
		values = params.GetFilterValues(gt)
		if len(values) == 0 {
			continue
		}
		switch gt {
		case MetaFilter:
			ids, err = v.Queries.ListLinkIdsByMetadata(ctx, sqlc.ListLinkIdsByMetadataParams{
				KvPairs:       values,
				LinksArchived: sqlc.NotArchived,
				LinksDeleted:  sqlc.NotDeleted,
			})
		case GroupTypeFilter:
			ids, err = v.Queries.ListLinkIdsByGroupType(ctx, sqlc.ListLinkIdsByGroupTypeParams{
				GroupTypes:    values,
				LinksArchived: sqlc.NotArchived,
				LinksDeleted:  sqlc.NotDeleted,
			})
		default:
			switch {
			case slices.Contains(values, "none"):
				ids, err = v.Queries.ListLinkIdsNotInGroupType(ctx, sqlc.ListLinkIdsNotInGroupTypeParams{
					GroupTypes:    []string{gt},
					LinksArchived: sqlc.NotArchived,
					LinksDeleted:  sqlc.NotDeleted,
				})
			default:
				ids, err = v.Queries.ListLinkIdsByGroupSlugs(ctx, sqlc.ListLinkIdsByGroupSlugsParams{
					Slugs:         values,
					LinksArchived: sqlc.NotArchived,
					LinksDeleted:  sqlc.NotDeleted,
				})
			}
		}
		if err != nil {
			goto end
		}
		if len(ids) == 0 {
			continue
		}
		// TODO: Once the UI supports calling API with multiple values this needs to be
		//       refactored to support AND logic vs. the OR logic it now has by default.
		linkIds = append(linkIds, ids...)
	}
	if len(linkIds) == 0 {
		html = safehtml.HTMLFromConstant("<div>No links for selection</div>")
		goto end
	} else {
		ll, err = v.Queries.ListFilteredLinks(ctx, sqlc.ListFilteredLinksParams{
			Ids:           linkIds,
			LinksArchived: sqlc.NotArchived,
			LinksDeleted:  sqlc.NotDeleted,
		})
	}
	if err != nil {
		goto end
	}
	links = linksFromLinkSet(ll)
	queryJSON, err = params.GetFilterJSON()
	if err != nil {
		slog.Error("Failed to get filter JSON", "err", err.Error())
		queryJSON = "{}"
	}
	html, err = linkSetTemplate.ExecuteToHTML(linkSet{
		apiURL:     makeURL(host),
		Links:      links,
		Label:      params.GetFilterLabels(),
		requestURI: requestURI,
		queryJSON:  queryJSON,
	})
	if err != nil {
		goto end
	}
end:
	return html, http.StatusInternalServerError, err
}

func linksFromLinkSet(ll []sqlc.ListFilteredLinksRow) (links []link) {
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
