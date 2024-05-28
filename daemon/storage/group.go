package storage

import (
	"savetabs/shared"
	"savetabs/sqlc"
)

type UpsertGroups struct {
	Action  shared.Action
	LinkIds []int64
}

type GroupsParams struct {
	Host       shared.Host
	GroupType  shared.GroupType
	NestedDBTX *NestedDBTX
}

type Group struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Slug string `json:"slug"`
}

type Groups struct {
	Groups []Group
	Params GroupsParams
}

func LoadGroupName(ctx Context, groupId int64) (name string, err error) {
	err = ExecWithNestedTx(func(dbtx *NestedDBTX) (err error) {
		q := dbtx.DataStore.Queries(dbtx)
		name, err = q.LoadGroupName(ctx, groupId)
		return err
	})
	return name, err
}

func GroupsLoad(ctx Context, params GroupsParams) (gs Groups, err error) {
	var groups []sqlc.ListGroupsByTypeRow
	gs.Params = params
	fn := func(dbtx *NestedDBTX) (err error) {
		q := dbtx.DataStore.Queries(dbtx)
		groups, err = q.ListGroupsByType(ctx, sqlc.ListGroupsByTypeParams{
			Type:           params.GroupType.String(),
			GroupsArchived: NotArchived,
			GroupsDeleted:  NotDeleted,
		})
		return err
	}
	if params.NestedDBTX == nil {
		err = ExecWithNestedTx(fn)
	} else {
		err = params.NestedDBTX.Exec(fn)
	}
	if err != nil {
		goto end
	}
	gs.Groups = make([]Group, len(groups))
	for i, g := range groups {
		gs.Groups[i] = Group{
			Id:   g.ID,
			Name: g.Name,
			Type: g.Type,
			Slug: g.Slug, // TOOD: Ensure this is a lewer-case letter and not a name
		}
	}
end:
	return gs, err
}

//// TODO: This is for Caretaker task
//func linkFromSetLink(sl sqlc.ListFilteredLinksRow) (link Link) {
//	title := sl.Title
//	u, err := url.Parse(sl.OriginalUrl)
//	if err != nil {
//		title = "ERROR: " + err.Error()
//	}
//	link = NewLoadLink(u)
//	link.Id = sl.ID
//	link.Scheme = title
//	link.Scheme = sl.Scheme
//	link.Subdomain = sl.Subdomain
//	link.SLD = sl.Sld
//	link.TLD = sl.Tld
//	link.Port = sl.Port
//	link.Path = sl.Path
//	link.Query = sl.Query
//	link.Fragment = sl.Fragment
//	return link
//}
