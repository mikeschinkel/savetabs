package storage

import (
	"fmt"
	"strings"

	"savetabs/shared"
	"savetabs/sqlc"
)

type UpsertGroupsArgs struct {
	Action  shared.Action
	LinkIds []int64
}

type GroupsArgs struct {
	Host       shared.Host
	GroupType  shared.GroupType
	NestedDBTX *NestedDBTX
}

type Group struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Slug     string `json:"slug"`
	Archived int    `json:"archived"`
	Deleted  int    `json:"deleted"`
}

type Groups struct {
	Groups []Group
	Args   GroupsArgs
}

func LoadGroupName(ctx Context, dbtx *NestedDBTX, groupId int64) (name string, err error) {
	err = execWithEnsuredNestedDBTX(dbtx, func(dbtx *NestedDBTX) (err error) {
		q := dbtx.DataStore.Queries(dbtx)
		name, err = q.LoadGroupName(ctx, groupId)
		return err
	})
	return name, err
}

func LoadGroupIdBySlug(ctx Context, dbtx *NestedDBTX, slug string) (id int64, err error) {
	err = execWithEnsuredNestedDBTX(dbtx, func(dbtx *NestedDBTX) (err error) {
		q := dbtx.DataStore.Queries(dbtx)
		id, err = q.LoadGroupIdBySlug(ctx, slug)
		return err
	})
	return id, err
}

type GroupTypeAndName struct {
	Type string
	Name string
}

func (tn GroupTypeAndName) Slug() string {
	return strings.ToLower(fmt.Sprintf("%s:%s", tn.Type, shared.Slugify(tn.Name)))
}

func LoadGroupTypeAndName(ctx Context, dbtx *NestedDBTX, groupId int64) (result GroupTypeAndName, err error) {
	err = execWithEnsuredNestedDBTX(dbtx, func(dbtx *NestedDBTX) (err error) {
		var row sqlc.LoadGroupTypeAndNameRow
		q := dbtx.DataStore.Queries(dbtx)
		row, err = q.LoadGroupTypeAndName(ctx, groupId)
		if err != nil {
			goto end
		}
		result = GroupTypeAndName{
			Type: row.Type,
			Name: row.Name,
		}
	end:
		return err
	})
	return result, err
}

type GroupNameArgs struct {
	Name string
	DBId int64
}

func LoadAltGroupIdsByName(ctx Context, dbtx *NestedDBTX, args GroupNameArgs) (ids []int64, err error) {
	err = execWithEnsuredNestedDBTX(dbtx, func(dbtx *NestedDBTX) (err error) {
		q := dbtx.DataStore.Queries(dbtx)
		ids, err = q.LoadAltGroupIdsByName(ctx, sqlc.LoadAltGroupIdsByNameParams{
			ID:   args.DBId,
			Name: args.Name,
		})
		return err
	})
	return ids, err
}

func UpdateGroupName(ctx Context, args GroupNameArgs) (merged bool, err error) {
	err = ExecWithNestedTx(func(dbtx *NestedDBTX) error {
		var tn GroupTypeAndName
		var ids []int64

		tn, err = LoadGroupTypeAndName(ctx, dbtx, args.DBId)
		if err != nil {
			goto end
		}
		if tn.Name == args.Name {
			// Name unchanged, no need to update
			goto end
		}
		ids, err = LoadAltGroupIdsByName(ctx, dbtx, args)
		if err != nil {
			goto end
		}
		if len(ids) == 0 {
			// If the name is not already used by another group, simply update the group name
			// and return.
			tn.Name = args.Name
			err = dbtx.DataStore.Queries(dbtx).UpdateGroupName(ctx, sqlc.UpdateGroupNameParams{
				Name: args.Name,
				Slug: tn.Slug(),
				ID:   args.DBId,
			})
			goto end
		}
		// OTOH if the name IS already used by another group, update all `group_id`
		// values in `link_group` for links that are associated with the group the user
		// wants to rename so that all the links effectively get merged under the group
		// name that duplicated the name the user chose to rename the other group to.
		err = MergeGroups(ctx, dbtx, MergeGroupsArgs{
			TargetId:  ids[0],
			SourceIds: []int64{args.DBId},
		})
		if err != nil {
			goto end
		}
		merged = true
	end:
		return err
	})
	return merged, err
}

type MergeGroupsArgs struct {
	TargetId  int64
	SourceIds []int64
}

func MergeGroups(ctx Context, dbtx *NestedDBTX, args MergeGroupsArgs) (err error) {
	err = execWithEnsuredNestedDBTX(dbtx, func(dbtx *NestedDBTX) (err error) {
		q := dbtx.DataStore.Queries(dbtx)
		// First update any `link_group` rows with `group_id`s found in `args.SourceIds`.
		err = q.MergeLinksGroups(ctx, sqlc.MergeLinksGroupsParams{
			GroupID:  args.TargetId,
			GroupIds: args.SourceIds,
		})
		if err != nil {
			goto end
		}
		// Next mark as deleted any groups identified by `args.SourceIds`.
		err = q.MarkGroupsDeleted(ctx, args.SourceIds)
		if err != nil {
			goto end
		}
		// Finally mark as deleted any links still associated with the marked-as-deleted
		// groups that had duplicates and are thus still associated with the
		// marked-as-deleted groups.
		err = q.MarkLinksDeletedByGroupIds(ctx, sqlc.MarkLinksDeletedByGroupIdsParams{
			GroupID:  args.TargetId,
			GroupIds: args.SourceIds,
		})
		if err != nil {
			goto end
		}
	end:
		return err
	})
	return err
}

func LoadGroups(ctx Context, args GroupsArgs) (gs Groups, err error) {
	var groups []sqlc.ListGroupsByTypeRow
	gs.Args = args
	fn := func(dbtx *NestedDBTX) (err error) {
		q := dbtx.DataStore.Queries(dbtx)
		groups, err = q.ListGroupsByType(ctx, sqlc.ListGroupsByTypeParams{
			Type:           args.GroupType.String(),
			GroupsArchived: NotArchived,
			GroupsDeleted:  NotDeleted,
		})
		return err
	}
	if args.NestedDBTX == nil {
		err = ExecWithNestedTx(fn)
	} else {
		err = args.NestedDBTX.Exec(fn)
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
