package ui

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/safehtml"
	"savetabs/shared"
	"savetabs/sqlc"
)

var groupsTemplate = GetTemplate("groups")

var _ MenuItemable = (*group)(nil)

type group struct {
	Id        int64
	Name      string
	Type      string
	TypeName  string
	LinkCount int64
	Links     []link
	Host      string
}

func (g group) LinksQueryParams() string {
	return fmt.Sprintf("g=%s", g.Slug())
}

func (g group) MenuItemType() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(GroupItemType)
}

func (g group) Slug() string {
	return strings.ToLower(g.Type) + "/" + shared.Slugify(g.Name)
}

func (g group) Url() string {
	return fmt.Sprintf("/html/groups/%s", g.Slug())
}

func (g group) Target() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`group-links`,
		strconv.FormatInt(g.Id, 10),
	)
}

func (g group) HTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`group`,
		strconv.FormatInt(g.Id, 10),
	)
}

func constructGroups(grs []sqlc.Group) []group {
	gg := make([]group, len(grs))
	for i, gr := range grs {
		gg[i] = newGroupFromSqlcGroup(gr)
	}
	return gg
}

func (v *Views) GetGroupHTML(ctx Context, host, gt, gs string) (html []byte, status int, err error) {
	var out bytes.Buffer

	//var gt groupType
	var ll []sqlc.ListLinksForGroupRow
	//var gs []group
	//
	ll, err = v.Queries.ListLinksForGroup(ctx, sqlc.ListLinksForGroupParams{
		GroupType: strings.ToUpper(gt),
		GroupSlug: gs,
	})
	if err != nil {
		goto end
	}
	if len(ll) == 0 {
		goto end
	}

	err = linksTemplate.Execute(&out, group{
		Host:      makeURL(host),
		Id:        ll[0].GroupID,
		Name:      ll[0].GroupName,
		Type:      ll[0].GroupType,
		TypeName:  ll[0].TypeName,
		LinkCount: int64(len(ll)),
		Links:     constructLinks(ll),
	})
	if err != nil {
		goto end
	}
	html = out.Bytes()
	goto end
end:
	return html, http.StatusInternalServerError, err
}

func constructLinks(ll []sqlc.ListLinksForGroupRow) []link {
	links := make([]link, len(ll))
	for i, rfg := range ll {
		r := &link{
			rowId: i + 1,
			Id:    rfg.ID.Int64,
			URL:   rfg.Url.String,
		}
		links[i] = *r
	}
	return links
}
