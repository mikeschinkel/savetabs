package model

import (
	"fmt"

	"savetabs/shared"
)

type MenuItem struct {
	Id      string
	LocalId string
	Menu    *Menu
	Label   string
}

type MenuItemParams struct {
	LocalId string
	Label   string
}

func newMenuItem(m *Menu, p MenuItemParams) MenuItem {
	return MenuItem{}.Renew(m, p)
}

var pristine MenuItem

func (mi MenuItem) Renew(m *Menu, p MenuItemParams) MenuItem {
	mi = pristine
	mi.Id = fmt.Sprintf("%s-%s", m.Type.String(), p.LocalId)
	mi.LocalId = p.LocalId
	mi.Menu = m
	mi.Label = p.Label
	return mi
}

type MenuItemLoadParams struct {
	MenuType shared.MenuType
	Menu     *Menu
}
type MenuItems struct {
	Items []MenuItem
}

func MenuItemsLoad(ctx Context, p MenuItemLoadParams) (items MenuItems, err error) {
	var groups Groups
	var gt shared.GroupType

	gt, err = shared.GroupTypeByCode(p.MenuType.String())
	if err != nil {
		// Panic because upstream should have cause this, so that needs to be where it is
		// fixed, not here. Hence failing here is a programming error.
		panic(err.Error())
	}

	groups, err = GroupsLoad(ctx, GroupsParams{
		GroupType: gt,
	})
	if err != nil {
		goto end
	}
	items.Items = make([]MenuItem, len(groups.Groups))
	for i, g := range groups.Groups {
		items.Items[i] = newMenuItem(p.Menu, MenuItemParams{
			LocalId: fmt.Sprintf("%s-%s", p.MenuType, g.Type.Lower()),
			Label:   g.Name,
		})
	}
end:
	return items, err
}
