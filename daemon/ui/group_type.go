package ui

import (
	"fmt"
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

// HTMLId uniquely identifies a Group Type across all entities that might
// appear in an HTML page.
func (gt groupType) HTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`gt`,
		strings.ToLower(gt.Type),
	)
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
