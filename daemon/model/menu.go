package model

import (
	"strings"

	"savetabs/shared"
	"savetabs/storage"
)

type Menu struct {
	Type  shared.MenuType
	Level int
	Items []MenuItem
}

func NewMenu(mt shared.MenuType, level int) Menu {
	return Menu{
		Type:  mt,
		Level: level,
		Items: make([]MenuItem, 0),
	}
}

type MenuParams struct {
	Type  shared.MenuType
	Level int
}

func MenuLoad(ctx Context, p MenuParams) (m Menu, err error) {
	var gts storage.GroupTypes
	var invalidGTWithStats storage.GroupType
	var cnt int

	m = NewMenu(p.Type, p.Level)
	err = storage.ExecWithNestedTx(func(dbtx *storage.NestedDBTX) (err error) {
		gts, err = storage.LoadGroupTypes(ctx, storage.GroupTypeParams{})
		if err != nil {
			goto end
		}
		invalidGTWithStats, err = storage.GroupTypeLoadWithStats(ctx, storage.GroupTypeParams{
			GroupType:  shared.GroupTypeInvalid,
			NestedDBTX: dbtx,
		})

		if err != nil {
			goto end
		}
	end:
		return err
	})
	if err != nil {
		goto end
	}
	m.Items = make([]MenuItem, len(gts.GroupTypes))
	cnt = len(m.Items)
	for i, gt := range gts.GroupTypes {
		if excludeGroupTypeAsMenuItem(gt, invalidGTWithStats) {
			cnt--
			continue
		}
		m.Items[i] = m.Items[i].Renew(&m, MenuItemParams{
			LocalId: strings.ToLower(gt.Type),
			Label:   gt.Plural,
		})
	}
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
