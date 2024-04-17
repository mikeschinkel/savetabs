package ui

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/safehtml"
	"savetabs/sqlc"
)

type IconState = safehtml.Identifier

var (
	BlankIcon     IconState = safehtml.IdentifierFromConstant("blank")
	ExpandedIcon  IconState = safehtml.IdentifierFromConstant("expanded")
	CollapsedIcon IconState = safehtml.IdentifierFromConstant("collapsed")
)

type menuItem struct {
	apiURL string
	Id     safehtml.Identifier
	Source MenuItemable
	Label  string
	menuItemArgs
}
type menuItemArgs struct {
	IconState    IconState
	DetailsClass string
	SummaryClass string
}

const topLevelSummaryClass = "py-4 my-0"

func newMenuItem(src MenuItemable, host, label string) menuItem {
	return newMenuItemWithArgs(src, host, label, menuItemArgs{
		SummaryClass: topLevelSummaryClass,
		IconState:    CollapsedIcon,
	})
}

func newMenuItemWithArgs(src MenuItemable, host, label string, args menuItemArgs) menuItem {
	mi := menuItem{
		apiURL:       makeURL(host),
		Source:       src,
		Label:        label,
		menuItemArgs: args,
	}
	mi.Id = mi.Identifier()
	return mi
}

func (mi menuItem) IconIsBlank() bool {
	return mi.IconState == BlankIcon
}

func (mi menuItem) Slug() safehtml.Identifier {
	return mi.Source.Identifier()
}

func (mi menuItem) LinksQueryParams() string {
	return "?" + mi.Source.LinksQueryParams()
}

func (mi menuItem) Identifier() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`mi`, mi.Slug().String())
}

func newGroupFromSqlcGroup(g sqlc.Group) group {
	return group{
		Id:       g.ID,
		Name:     g.Name,
		Type:     g.Type,
		TypeName: g.Name,
	}
}

func menuItemsFromListGroupTypesRows(host string, gtrs []sqlc.ListGroupsTypeRow) []menuItem {
	var menuItems []menuItem

	cnt := len(gtrs)

	// No need to show invalid as a group type if
	// there are no resources of that type
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
	menuItems = make([]menuItem, cnt)
	for i, gtr := range gtrs {
		if i == invalid {
			continue
		}
		src := newGroupTypeFromListGroupsTypeRow(gtr)
		menuItems[i] = newMenuItem(src, host, gtr.Plural.String)
	}
	menuItems = append(menuItems,
		newMenuItemWithArgs(allLinks{}, host, "All Links", menuItemArgs{
			SummaryClass: topLevelSummaryClass,
			IconState:    BlankIcon,
		}),
	)
	return menuItems
}

var _ MenuItemable = (*allLinks)(nil)

type allLinks struct{}

func (allLinks) LinksQueryParams() string {
	return "all=1"
}

func (a allLinks) Identifier() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`gt-all`)
}

func (a allLinks) MenuItemType() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`A`)
}

func (v *Views) GetMenuItemHTML(ctx Context, host, item string) (html []byte, status int, err error) {
	var out bytes.Buffer
	var items []menuItem

	items, err = v.getMenuItemsForType(ctx, host, item)
	if err != nil {
		goto end
	}
	err = menuTemplate.Execute(&out, menu{
		apiURL:    makeURL(host),
		MenuItems: items,
	})
	if err != nil {
		goto end
	}
	html = out.Bytes()
end:
	return html, http.StatusInternalServerError, err
}

type ItemType string

const (
	GroupTypeItemType = "gt"
	GroupItemType     = "grp"
)

var matchMenuItemKey = regexp.MustCompile(`^(gt|grp)-(.+)$`)

func (v *Views) getMenuItemsForType(ctx Context, host, key string) (items []menuItem, err error) {
	var keys []string
	var gt sqlc.GroupType

	if !matchMenuItemKey.MatchString(key) {
		err = errors.Join(ErrInvalidKeyFormat, fmt.Errorf(`key=%s`, key))
		goto end
	}
	keys = matchMenuItemKey.FindStringSubmatch(key)
	switch keys[1] {
	case GroupTypeItemType: // Group Type
		var gs []sqlc.Group
		gs, err = v.Queries.ListGroupsByType(ctx, strings.ToUpper(keys[2]))
		if err != nil {
			goto end
		}
		gt, err = v.Queries.LoadGroupType(ctx, strings.ToUpper(keys[2]))
		if err != nil {
			goto end
		}
		err = nil
		items = menuItemsFromGroups(host, gt.Plural.String, gs)
	}
end:
	return items, err
}

func menuItemsFromGroups(host, gt string, gs []sqlc.Group) []menuItem {
	var menuItems []menuItem
	args := menuItemArgs{
		IconState:    BlankIcon,
		DetailsClass: "p-0 m-0",
		SummaryClass: "p-0 m-0",
	}

	menuItems = make([]menuItem, len(gs)+1)
	menuItems[0] = newMenuItemWithArgs(noMenuItem{}, host, fmt.Sprintf("<No %s>", gt), args)
	for i, g := range gs {
		src := newGroupFromSqlcGroup(g)
		menuItems[i+1] = newMenuItemWithArgs(src, host, g.Name, args)
	}
	return menuItems
}

var _ MenuItemable = (*noMenuItem)(nil)

type noMenuItem struct{}

func (noMenuItem) LinksQueryParams() string {
	return "g=none"
}

func (i noMenuItem) MenuItemType() safehtml.Identifier {
	return safehtml.IdentifierFromConstant("_")
}
func (noMenuItem) Identifier() safehtml.Identifier {
	return safehtml.IdentifierFromConstant("none")
}
