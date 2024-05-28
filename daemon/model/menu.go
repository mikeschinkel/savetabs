package model

import (
	"strings"

	"savetabs/shared"
	"savetabs/storage"
)

type Menu struct {
	Type  *shared.MenuType
	Items []MenuItem
}

func NewMenu(mt *shared.MenuType) *Menu {
	return &Menu{
		Type:  mt,
		Items: make([]MenuItem, 0),
	}
}

type MenuParams struct {
	Type  *shared.MenuType
	Level int
}

// LoadMenu loads the top level menu
func LoadMenu(ctx Context, p MenuParams) (m *Menu, err error) {
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
		mt, err := shared.MenuTypeByParentTypeAndMenuName(shared.GroupTypeMenuType, gt.Type)
		if err != nil {
			shared.Panicf("Invalid group type '%s' loaded from database", gt.Type)
		}
		cnt++
		item = item.Renew(MenuItemArgs{
			LocalId:  strings.ToLower(gt.Type),
			Label:    gt.Plural,
			Menu:     m,
			MenuType: mt,
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
