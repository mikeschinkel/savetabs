package ui

import (
	"bytes"
	"fmt"

	"github.com/google/safehtml"
	"savetabs/shared"
	"savetabs/sqlc"
)

var groupsTemplate = getTemplate("groups")

type group struct {
	Id            int64
	Name          string
	Type          string
	TypeName      string
	ResourceCount int64
}

func (g group) Slug() string {
	return shared.Slugify(g.Name)
}

func (g group) Url() string {
	return fmt.Sprintf("/html/groups/%s", g.Slug())
}

func (g group) Target() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`group-resources`, g.Slug())
}

func GroupsByGroupTypeHTML(ctx Context, host, typeName string) (html []byte, err error) {
	var gt groupType
	var out bytes.Buffer
	var gg []sqlc.GroupsWithCount
	var gs []group

	t, err := groupTypeFromName(ctx, typeName)
	if err != nil {
		goto end
	}
	gg, err = queries.ListGroupsWithCountsByGroupType(ctx, t)
	if err != nil {
		goto end
	}
	gs = constructGroups(gg)
	gt = groupType{
		Host:   "http://" + host,
		Groups: gs,
	}
	err = groupsTemplate.Execute(&out, gt)
	if err != nil {
		goto end
	}
	html = out.Bytes()
end:
	return html, err
}

func constructGroups(grs []sqlc.GroupsWithCount) []group {
	gg := make([]group, len(grs))
	for i, gr := range grs {
		g := &group{
			Id:            gr.ID,
			Name:          gr.Name.String,
			Type:          gr.Type.(string),
			TypeName:      gr.TypeName.String,
			ResourceCount: gr.ResourceCount,
		}
		gg[i] = *g
	}
	return gg
}
