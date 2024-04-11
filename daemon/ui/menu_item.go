package ui

import (
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
	Icon   IconType
}

func newMenuItem(src MenuItemable, host, label string) menuItem {
	return newMenuItemWithIcon(src, host, label, ExpandIcon)
}

func newMenuItemWithIcon(src MenuItemable, host, label string, icon IconType) menuItem {
	if icon == "" {
		icon = "blank"
	}
	mi := menuItem{
		apiURL: makeURL(host),
		Source: src,
		Label:  label,
		Icon:   icon,
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
		newMenuItemWithIcon(allLinks{}, host, "All Links", BlankIcon),
	)
	return menuItems
}

var _ MenuItemable = (*allLinks)(nil)

type allLinks struct{}

func (a allLinks) Identifier() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`gt-all`)
}
