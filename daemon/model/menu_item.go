package model

import (
	"github.com/mikeschinkel/savetabs/daemon/shared"
)

type MenuItem struct {
	*Menu
	LocalId     string
	Label       string
	FilterType  *shared.FilterType
	ContextMenu *shared.ContextMenu
	DBId        int64
}

type MenuItemArgs struct {
	LocalId     string
	Label       string
	MenuType    *shared.MenuType
	Menu        *Menu
	ContextMenu *shared.ContextMenu
	DBId        int64
}

func newMenuItem(p MenuItemArgs) MenuItem {
	return MenuItem{}.Renew(p)
}

var zeroStateMenuItem MenuItem

func (mi MenuItem) Renew(args MenuItemArgs) MenuItem {
	mi = zeroStateMenuItem
	if args.LocalId == "" {
		args.LocalId = shared.Slugify(args.Label)
	}
	mi.DBId = args.DBId
	mi.LocalId = args.LocalId
	mi.Menu = args.Menu
	mi.Label = args.Label
	mi.ContextMenu = args.ContextMenu

	if args.MenuType == nil {
		mt, err := shared.MenuTypeByParentTypeAndMenuName(args.Menu.Type, args.LocalId)
		if err != nil {
			panic(err.Error())
		}
		args.MenuType = mt
	}
	return mi
}

type LoadMenuItemParams struct {
	MenuType *shared.MenuType
}
type MenuItems struct {
	Items []MenuItem
}

func LoadMenuItems(ctx Context, p LoadMenuItemParams) (items MenuItems, err error) {
	var groups Groups
	var gt shared.GroupType
	var menu *Menu

	gt, err = shared.ParseGroupTypeByLetter(p.MenuType.Name())
	if err != nil {
		// Panic because upstream should have cause this, so that needs to be where it is
		// fixed, not here. Hence failing here is a programming error.
		panic(err.Error())
	}

	groups, err = LoadGroups(ctx, GroupsArgs{
		GroupType: gt,
	})
	if err != nil {
		goto end
	}

	menu = NewMenu(p.MenuType)

	items.Items = shared.ConvertSlice(groups.Groups, func(grp Group) MenuItem {
		return newMenuItem(MenuItemArgs{
			LocalId:     shared.Slugify(grp.Name),
			Label:       grp.Name,
			Menu:        menu,
			MenuType:    shared.GroupTypeMenuType,
			ContextMenu: shared.NewContextMenu(shared.GroupContextMenuType, grp.Id),
			DBId:        grp.Id,
		})
	})
end:
	return items, err
}
