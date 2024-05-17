package model

import (
	"fmt"

	"savetabs/shared"
	"savetabs/storage"
)

type Group struct {
	Id       int64
	Name     string
	Type     shared.GroupType
	TypeName string
}

func (grp Group) Slug() string {
	return fmt.Sprintf("%s/%s", grp.Type.Lower(), shared.Slugify(grp.Name))
}

type Groups struct {
	Groups []Group
	Params GroupsParams
}

type GroupsParams storage.GroupsParams

func NewGroups(groups storage.Groups) Groups {
	gs := make([]Group, len(groups.Groups))
	for i, grp := range groups.Groups {
		gt, err := shared.GroupTypeByCode(grp.Type)
		if err != nil {
			// Panic because upstream should have cause this, so that needs to be where it is
			// fixed, not here. Hence failing here is a programming error.
			panic(err.Error())
		}
		gs[i] = Group{
			Id:       grp.Id,
			Name:     grp.Name,
			Type:     gt,
			TypeName: grp.TypeName,
		}
	}
	return Groups{
		Groups: gs,
		Params: GroupsParams(groups.Params),
	}
}

func GroupsLoad(ctx Context, params GroupsParams) (groups Groups, err error) {
	var gs storage.Groups
	gs, err = storage.GroupsLoad(ctx, storage.GroupsParams(params))
	if err != nil {
		goto end
	}
end:
	return NewGroups(gs), err
}
