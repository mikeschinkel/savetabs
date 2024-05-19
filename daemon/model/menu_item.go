package model

import (
	"fmt"

	"github.com/google/safehtml"
	"savetabs/shared"
)

var _ shared.MenuItemable = (*MenuItem)(nil)

type MenuItem struct {
	shared.MenuItemable
	LocalId string
	Label   string
	Type    *shared.MenuType
}

func (mi MenuItem) Root() shared.MenuItemable {
	parent := mi.MenuItemable.Parent()
	for parent != nil {
		parent = parent.Parent()
	}
	return parent
}

func (mi MenuItem) Level() int {
	return mi.MenuItemable.Level()
}

func (mi MenuItem) MenuType() *shared.MenuType {
	return mi.Type
}

func (mi MenuItem) HTMLId() safehtml.Identifier {
	return shared.MakeSafeId(fmt.Sprintf("%s-%s",
		mi.MenuItemable.HTMLId(),
		mi.LocalId,
	))
}

func (mi MenuItem) Parent() shared.MenuItemable {
	return mi.MenuItemable
}

func (mi MenuItem) SubmenuURL() safehtml.URL {
	return shared.MakeSafeURL(fmt.Sprintf("%s/%s",
		mi.MenuItemable.SubmenuURL(), // e.g.	menu/gt/g
		mi.LocalId,                   // e.g. `golang`
	))
}

func (mi MenuItem) ItemURL() safehtml.URL {

	m := mi.MenuItemable
	return shared.MakeSafeURL(fmt.Sprintf("%s%s",
		m.ItemURL(), // e.g. `linkset/?gt=g`
		mi.MenuType().Params(),
	))
}

type MenuItemParams struct {
	LocalId string
	Label   string
	Menu    shared.MenuItemable
	Type    *shared.MenuType
}

func newMenuItem(p MenuItemParams) MenuItem {
	return MenuItem{}.Renew(p)
}

var pristineMenuItem MenuItem

func (mi MenuItem) Renew(p MenuItemParams) MenuItem {
	mi = pristineMenuItem
	if p.LocalId == "" {
		p.LocalId = shared.Slugify(p.Label)
	}
	mi.LocalId = p.LocalId
	mi.MenuItemable = p.Menu
	mi.Label = p.Label

	if p.Type == nil {
		mt, err := shared.MenuTypeByParentTypeAndMenuName(p.Menu.MenuType(), p.LocalId)
		if err != nil {
			panic(err.Error())
		}
		p.Type = &mt
	}
	mi.Type = p.Type
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
		return newMenuItem(MenuItemParams{
			LocalId: gt.String(),
			Label:   grp.Name,
			Menu:    p.Menu,
			Type:    &shared.GroupTypeMenuType,
		})
	})
end:
	return items, err
}
