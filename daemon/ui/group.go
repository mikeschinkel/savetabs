package ui

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/safehtml"
	"savetabs/shared"
	"savetabs/sqlc"
)

var groupsTemplate = GetTemplate("groups")

type group struct {
	Id            int64
	Name          string
	Type          string
	TypeName      string
	ResourceCount int64
	Resources     []resource
	Host          string
}

func (g group) Slug() string {
	return strings.ToLower(g.Type) + "/" + shared.Slugify(g.Name)
}

func (g group) Url() string {
	return fmt.Sprintf("/html/groups/%s", g.Slug())
}

func (g group) Target() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`group-resources`,
		strconv.FormatInt(g.Id, 10),
	)
}

func (g group) Identifier() safehtml.Identifier {
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

func GetGroupHTML(ctx Context, host, gt, gs string) (html []byte, err error) {
	var out bytes.Buffer

	//var gt groupType
	var rfgs []sqlc.ListResourcesForGroupRow
	//var gs []group
	//
	rfgs, err = queries.ListResourcesForGroup(ctx, sqlc.ListResourcesForGroupParams{
		GroupType: strings.ToUpper(gt),
		GroupSlug: gs,
	})
	if err != nil {
		goto end
	}
	if len(rfgs) == 0 {
		goto end
	}

	err = resourcesTemplate.Execute(&out, group{
		Host:          makeURL(host),
		Id:            rfgs[0].GroupID,
		Name:          rfgs[0].GroupName,
		Type:          rfgs[0].GroupType,
		TypeName:      rfgs[0].TypeName,
		ResourceCount: int64(len(rfgs)),
		Resources:     constructResources(rfgs),
	})
	if err != nil {
		goto end
	}
	html = out.Bytes()
	goto end
end:
	return html, err
}
