package ui

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/safehtml"
	"savetabs/sqlc"
)

type IconType string

const (
	BlankIcon    IconType = "blank"
	ExpandIcon   IconType = "right-chevron"
	CollapseIcon IconType = "down-chevron"
)

type menuItem struct {
	apiURL string
	Id     safehtml.Identifier
	Source MenuItemable
	Label  string
	menuItemArgs
}
type menuItemArgs struct {
	Icon         IconType
	DetailsClass string
	SummaryClass string
}

func newMenuItem(src MenuItemable, host, label string) menuItem {
	return newMenuItemWithArgs(src, host, label, menuItemArgs{
		SummaryClass: "py-4 my-0",
		Icon:         ExpandIcon,
	})
}

func newMenuItemWithArgs(src MenuItemable, host, label string, args menuItemArgs) menuItem {
	if args.Icon == "" {
		args.Icon = "blank"
	}
	mi := menuItem{
		apiURL:       makeURL(host),
		Source:       src,
		Label:        label,
		menuItemArgs: args,
	}
	mi.Id = mi.Identifier()
	return mi
}

func (mi menuItem) Slug() safehtml.Identifier {
	return mi.Source.Identifier()
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
		if gtr.ResourceCount != 0 {
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
			Icon: BlankIcon,
		}),
	)
	return menuItems
}

var _ MenuItemable = (*allLinks)(nil)

type allLinks struct{}

func (a allLinks) Identifier() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`gt-all`)
}

func MenuItemHTML(host string, item string) (html []byte, err error) {
	var out bytes.Buffer
	var items []menuItem

	items, err = GetMenuItemsForType(context.Background(), host, item)
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
	return html, err
}

type ItemType string

const (
	GroupTypeItemType = "gt"
)

func GetMenuItemsForType(ctx Context, host, key string) (items []menuItem, err error) {
	keys := strings.SplitAfterN(key, "-", 2)
	if len(keys) != 2 {
		err = errors.Join(ErrInvalidKeyFormat, fmt.Errorf(`key=%s`, key))
		goto end
	}
	switch strings.TrimRight(keys[0], "-") {
	case GroupTypeItemType: // Group Type
		var gs []sqlc.Group
		gs, err = queries.ListGroupsByType(ctx, strings.ToUpper(keys[1]))
		if err != nil {
			goto end
		}
		items = menuItemsFromGroups(host, gs)
	}
end:
	return items, err
}

func menuItemsFromGroups(host string, gs []sqlc.Group) []menuItem {
	var menuItems []menuItem

	menuItems = make([]menuItem, len(gs))
	for i, g := range gs {
		src := newGroupFromSqlcGroup(g)
		menuItems[i] = newMenuItemWithArgs(src, host, g.Name, menuItemArgs{
			Icon:         ExpandIcon,
			DetailsClass: "p-0 m-0",
			SummaryClass: "px-0 py-8 m-0",
		})
	}
	return menuItems
}
