package model

import (
	"strings"

	"github.com/google/safehtml"
	"savetabs/shared"
	"savetabs/storage"
)

var _ shared.MenuItemable = (*Menu)(nil)

type Menu struct {
	Type   *shared.MenuType
	Items  []MenuItem
	level  int
	parent *Menu
}

func (m Menu) SubmenuURL() (u safehtml.URL) {
	if m.level == 0 {
		u = shared.MakeSafeURL("/")
		goto end
	}
	u = shared.MakeSafeURLf("/%s", strings.Join(m.Type.Slice(), "/"))
end:
	return u
}

func (m Menu) Parent() shared.MenuItemable {
	return m.parent
}

func (m Menu) ItemURL() (u safehtml.URL) {
	return shared.MakeSafeURL("?")
}

func (m Menu) Level() int {
	return m.level
}

func (m Menu) MenuType() *shared.MenuType {
	return m.Type
}

func (m Menu) HTMLId() (id safehtml.Identifier) {
	if m.level == 0 {
		id = safehtml.IdentifierFromConstant("mi")
		goto end
	}
	id = shared.MakeSafeIdf("mi-%s", strings.Join(m.Type.Slice(), "-"))
end:
	return id
}

func NewMenu(mt *shared.MenuType) *Menu {
	var parent *Menu
	level := mt.Level()
	if level > 0 {
		parent = NewMenu(mt.Parent)
	}
	return &Menu{
		Type:   mt,
		Items:  make([]MenuItem, 0),
		level:  level,
		parent: parent,
	}
}

type MenuParams struct {
	Type  *shared.MenuType
	Level int
}

// MenuLoad loads the top level menu
func MenuLoad(ctx Context, p MenuParams) (m *Menu, err error) {
	var gts storage.GroupTypes
	var invalidGTWithStats storage.GroupType
	var cnt int

	m = NewMenu(p.Type)
	err = storage.ExecWithNestedTx(func(dbtx *storage.NestedDBTX) (err error) {
		gts, err = storage.LoadGroupTypes(ctx, storage.GroupTypeParams{})
		if err != nil {
			goto end
		}
		invalidGTWithStats, err = storage.GroupTypeLoadWithStats(ctx, storage.GroupTypeParams{
			GroupType:  shared.GroupTypeInvalid,
			NestedDBTX: dbtx,
		})
	end:
		return err
	})
	if err != nil {
		goto end
	}
	m.Items = shared.ConvertSliceWithFilter(gts.GroupTypes, func(gt storage.GroupType) (item MenuItem, _ bool) {
		if excludeGroupTypeAsMenuItem(gt, invalidGTWithStats) {
			return item, false
		}
		typ, err := shared.MenuTypeByValue(gt.Type)
		if err != nil {
			shared.Panicf("Invalid group type '%s' loaded from database", gt.Type)
		}
		cnt++
		item = m.Items[cnt-1].Renew(MenuItemParams{
			LocalId: strings.ToLower(gt.Type),
			Label:   gt.Plural,
			Menu:    m,
			Type:    &typ,
		})
		return item, true
	})
	if cnt < len(m.Items) {
		m.Items = m.Items[:cnt]
	}
end:
	return m, err
}

var invalidGroupType = shared.GroupTypeInvalid.String()

// includeGroupTypeAsMenuItem will return true unless the group type is `Invalid`
// and it has no active links, e.g.` active` means not deleted and not archived.)
func excludeGroupTypeAsMenuItem(gt storage.GroupType, invalidGTWithStats storage.GroupType) (exclude bool) {
	if gt.Type != invalidGroupType {
		goto end
	}
	if invalidGTWithStats.HasActiveLinks() {
		goto end
	}
	exclude = true
end:
	return exclude
}
