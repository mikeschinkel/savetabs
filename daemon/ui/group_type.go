package ui

import (
	"bytes"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/google/safehtml"
	"savetabs/shared"
	"savetabs/sqlc"
)

var _ MenuItemable = (*groupType)(nil)

type groupType struct {
	Type       string
	Name       string
	Plural     string
	GroupCount int64
	Groups     []group
	Host       string
	Sort       int8
}

func (gt groupType) LinksQueryParams() string {
	return fmt.Sprintf("gt=%s", strings.ToLower(gt.Type))
}

func (gt groupType) MenuItemType() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(GroupTypeItemType)
}

// Slug uniquely identifies a Group Type
func (gt groupType) Slug() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`gt`,
		shared.Slugify(gt.Name),
	)
}

func (gt groupType) Url() string {
	return fmt.Sprintf(`/html/group-types/%s/groups`, gt.Slug())
}

func (gt groupType) Target() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`gtg`,
		strings.ToLower(gt.Type))
}

// Identifier uniquely identifies a Group Type across all entities that might
// appear in an HTML page.
func (gt groupType) Identifier() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`gt`,
		strings.ToLower(gt.Type),
	)
}

func getGroupTypeMap(ctx Context) (gtm groupTypeMap, err error) {
	var gts []sqlc.ListGroupsTypeRow

	gts, err = queries.ListGroupsType(ctx)
	if err != nil {
		goto end
	}
	gtm = newGroupTypeMap(gts)
end:
	return gtm, err
}

func groupTypeFromName(ctx Context, name string) (t string, err error) {
	var ok bool
	var gtm groupTypeMap
	var gt groupType

	gtm, err = getGroupTypeMap(ctx)
	if err != nil {
		goto end
	}
	gt, ok = gtm[name]
	if !ok {
		t = "I"
		goto end
	}
	t = gt.Type
end:
	return t, err
}

type groupTypeMap map[string]groupType

func (gtm groupTypeMap) Map(fn func(gt groupType) groupType) groupTypeMap {
	for i, x := range gtm {
		gtm[i] = fn(x)
	}
	return gtm
}

func (gtm groupTypeMap) AsSortedSlice() (gts []groupType) {
	gts = make([]groupType, len(gtm))
	i := 0
	for _, gt := range gtm {
		gts[i] = gt
		i++
	}
	slices.SortFunc(gts, func(a, b groupType) int {
		switch {
		case a.Sort > b.Sort:
			return 1
		case a.Sort < b.Sort:
			return -1
		}
		return 0
	})
	return gts
}

func newGroupTypeFromListGroupsTypeRow(gtr sqlc.ListGroupsTypeRow) groupType {
	return groupType{
		Type:       gtr.Type,
		Name:       gtr.Name.String,
		Plural:     gtr.Plural.String,
		GroupCount: gtr.GroupCount,
		Sort:       int8(gtr.Sort.Int64),
	}
}

func newGroupTypeMap(gtrs []sqlc.ListGroupsTypeRow) groupTypeMap {
	cnt := len(gtrs)

	// No need to show invalid as a group type if
	// there are no links of that type
	invalid := -1
	for i, gtr := range gtrs {
		if gtr.LinkCount != 0 {
			continue
		}
		if gtr.Type != "I" {
			continue
		}
		cnt--
		invalid = i
		break
	}
	gts := make(groupTypeMap, cnt)
	for i, gtr := range gtrs {
		if i == invalid {
			continue
		}
		gts[strings.ToLower(gtr.Name.String)] = newGroupTypeFromListGroupsTypeRow(gtr)
	}
	return gts
}

func GetGroupTypeGroupsHTML(ctx Context, host, groupTypeName string) (html []byte, status int, err error) {
	var gt groupType
	var out bytes.Buffer
	var gg []sqlc.Group
	var gs []group

	t, err := groupTypeFromName(ctx, groupTypeName)
	if err != nil {
		goto end
	}
	gg, err = queries.ListGroupsByType(ctx, strings.ToUpper(t))
	if err != nil {
		goto end
	}
	gs = constructGroups(gg)
	gt = groupType{
		Host:   makeURL(host),
		Groups: gs,
	}
	err = groupsTemplate.Execute(&out, gt)
	if err != nil {
		goto end
	}
	html = out.Bytes()
end:
	return html, http.StatusInternalServerError, err
}
