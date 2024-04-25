package ui

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/google/safehtml"
	"savetabs/sqlc"
)

type linkSet struct {
	apiURL   string
	rawQuery string
	Label    string
	Links    []link
}

func (ls linkSet) HeaderHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`links-row-head`)
}
func (ls linkSet) FooterHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`links-row-foot`)
}
func (ls linkSet) URLQuery() safehtml.URL {
	return safehtml.URLSanitized("?" + ls.rawQuery)
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

func (v *Views) GetLinkSetHTML(ctx Context, host string, params FilterValueGetter, rawQuery string) (html []byte, status int, err error) {
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
			switch {
			case slices.Contains(values, "none"):
				ids, err = v.Queries.ListLinkIdsNotInGroupType(ctx, sqlc.ListLinkIdsNotInGroupTypeParams{
					GroupTypes:   []string{gt},
					LinkArchived: sqlc.NotArchived,
				})
			default:
				ids, err = v.Queries.ListLinkIdsByGroupSlugs(ctx, sqlc.ListLinkIdsByGroupSlugsParams{
					Slugs:        values,
					LinkArchived: sqlc.NotArchived,
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
	links = linksFromLinkSet(ll)
	err = linkSetTemplate.Execute(&out, linkSet{
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
