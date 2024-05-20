package model

import (
	"savetabs/shared"
)

type MenuItem struct {
	*Menu
	LocalId string
	Label   string
	Type    *shared.MenuType
}

type MenuItemArgs struct {
	LocalId string
	Label   string
	Type    *shared.MenuType
	Menu    *Menu
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
	mi.LocalId = args.LocalId
	mi.Menu = args.Menu
	mi.Label = args.Label

	if args.Type == nil {
		mt, err := shared.MenuTypeByParentTypeAndMenuName(args.Menu.Type, args.LocalId)
		if err != nil {
			panic(err.Error())
		}
		args.Type = mt
	}
	mi.Type = args.Type
	return mi
}

type LoadMenuItemParams struct {
	MenuType *shared.MenuType
	Menu     *Menu
}
type MenuItems struct {
	Items []MenuItem
}

func LoadMenuItems(ctx Context, p LoadMenuItemParams) (items MenuItems, err error) {
	var groups Groups
	var gt shared.GroupType

	gt, err = shared.GroupTypeByType(p.MenuType.Name())
	if err != nil {
		// Panic because upstream should have cause this, so that needs to be where it is
		// fixed, not here. Hence failing here is a programming error.
		panic(err.Error())
	}

	groups, err = LoadGroups(ctx, GroupsParams{
		GroupType: gt,
	})
	if err != nil {
		goto end
	}

	items.Items = shared.ConvertSlice(groups.Groups, func(grp Group) MenuItem {
		return newMenuItem(MenuItemArgs{
			LocalId: gt.String(),
			Label:   grp.Name,
			Menu:    p.Menu,
			Type:    shared.GroupTypeMenuType,
		})
	})
end:
	return items, err
}
