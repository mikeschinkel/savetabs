package storage

import (
	"savetabs/shared"
	"savetabs/sqlc"
)

type GroupTypes struct {
	GroupTypes []GroupType
}

type GroupType struct {
	Type           string `json:"type"`
	Name           string `json:"name"`
	Plural         string `json:"plural"`
	GroupCount     int    `json:"group_count"`
	GroupsArchived int    `json:"groups_archived"`
	GroupsDeleted  int    `json:"groups_deleted"`
	LinkCount      int    `json:"link_count"`
	LinksArchived  int    `json:"links_archived"`
	LinksDeleted   int    `json:"links_deleted"`
	Sort           int    `json:"sort"`
}

func (gt GroupType) HasActiveLinks() bool {
	return gt.LinkCount-gt.LinksArchived-0-gt.LinksDeleted > 0
}

type GroupTypeParams struct {
	// GroupType used for queries that are specific to group type, but not when not.
	GroupType shared.GroupType

	// NestedDBTX only used when caller is already participating in a transaction.
	NestedDBTX *NestedDBTX
}

func execWithEnsuredNestedDBTX(dbtx *NestedDBTX, fn func(dbtx *NestedDBTX) error) (err error) {
	if dbtx != nil {
		err = dbtx.Exec(fn)
		goto end
	}
	err = ExecWithNestedTx(fn)
end:
	return err
}

// GroupTypesLoad loads a full list of group types from the `group_type` table of the data store.
// It includes link stats for each group.
func LoadGroupTypes(ctx Context, p GroupTypeParams) (gts GroupTypes, err error) {
	var rows []sqlc.GroupsType
	err = execWithEnsuredNestedDBTX(p.NestedDBTX, func(dbtx *NestedDBTX) (err error) {
		q := dbtx.DataStore.Queries(dbtx)
		rows, err = q.ListGroupsType(ctx)
		return err
	})
	if err != nil {
		goto end
	}
	gts.GroupTypes = make([]GroupType, len(rows))
	for i, row := range rows {
		gts.GroupTypes[i] = GroupType{
			Type:           row.Type,
			Name:           row.Name.String,
			Plural:         row.Plural.String,
			GroupCount:     int(row.GroupCount),
			GroupsArchived: int(row.GroupsArchived),
			GroupsDeleted:  int(row.GroupsDeleted),
			LinkCount:      int(row.LinkCount),
			LinksArchived:  int(row.LinksArchived),
			LinksDeleted:   int(row.LinksDeleted),
			Sort:           int(row.Sort.Int64),
		}
	}
end:
	return gts, err
}

type GroupsType sqlc.GroupsType

func GroupTypeLoadWithStats(ctx Context, p GroupTypeParams) (gt GroupType, err error) {
	var row sqlc.GroupsType
	err = execWithEnsuredNestedDBTX(p.NestedDBTX, func(dbtx *NestedDBTX) (err error) {
		q := dbtx.DataStore.Queries(dbtx)
		row, err = q.LoadGroupTypeWithStats(ctx, p.GroupType.String())
		return err
	})
	if err != nil {
		goto end
	}
	gt = GroupType{
		Type:           row.Type,
		Name:           row.Name.String,
		Plural:         row.Plural.String,
		GroupCount:     int(row.GroupCount),
		GroupsArchived: int(row.GroupsArchived),
		GroupsDeleted:  int(row.GroupsDeleted),
		LinkCount:      int(row.LinkCount),
		LinksArchived:  int(row.LinksArchived),
		LinksDeleted:   int(row.LinksDeleted),
		Sort:           int(row.Sort.Int64),
	}
end:
	return gt, err
}
