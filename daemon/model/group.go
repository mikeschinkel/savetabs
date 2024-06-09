package model

import (
	"fmt"

	"savetabs/shared"
	"savetabs/storage"
)

type Group struct {
	Id   int64
	Name string
	Type shared.GroupType
}

func (grp Group) Slug() string {
	return fmt.Sprintf("%s:%s", grp.Type.Lower(), shared.Slugify(grp.Name))
}

type Groups struct {
	GroupsArgs
	Groups []Group
}

type GroupsArgs storage.GroupsArgs

func NewGroups(groups storage.Groups) Groups {
	gs := make([]Group, len(groups.Groups))
	for i, grp := range groups.Groups {
		gt, err := shared.ParseGroupTypeByLetter(grp.Type)
		if err != nil {
			// Panic because upstream should have cause this, so that needs to be where it is
			// fixed, not here. Hence failing here is a programming error.
			panic(err.Error())
		}
		gs[i] = Group{
			Id:   grp.Id,
			Name: grp.Name,
			Type: gt,
		}
	}
	return Groups{
		GroupsArgs: GroupsArgs(groups.Args),
		Groups:     gs,
	}
}

func LoadGroupName(ctx Context, groupId int64) (name string, err error) {
	return storage.LoadGroupName(ctx, nil, groupId)
}

func LoadGroupIdBySlug(ctx Context, slug string) (int64, error) {
	return storage.LoadGroupIdBySlug(ctx, nil, slug)
}

func LoadGroups(ctx Context, params GroupsArgs) (groups Groups, err error) {
	var gs storage.Groups
	gs, err = storage.LoadGroups(ctx, storage.GroupsArgs(params))
	if err != nil {
		goto end
	}
end:
	return NewGroups(gs), err
}
