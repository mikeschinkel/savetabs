package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/safehtml"
	"savetabs/shared"
)

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
