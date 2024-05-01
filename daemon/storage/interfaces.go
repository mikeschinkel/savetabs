package storage

import (
	"savetabs/shared"
)

type Modifier interface {
	UpsertLinksWithGroups(Context, LinksWithGroups) error
}

type LinksWithGroupsGetSetter interface {
	GetLinkCount() int
	GetLinksWithGroups() []LinkWithGroupGetSetter
}
type LinkWithGroupGetSetter interface {
	LinkGetSetter
	GetGroup() string
	SetGroup(string)
	GetGroupId() int64
	SetGroupId(int64)
	GetGroupType() string
	SetGroupType(string)
}

type LinkSetGetSetter interface {
	GetLinkCount() int
	GetLinks() []LinkGetSetter
}
type LinkGetSetter interface {
	GetId() int64
	SetId(int64)
	GetTitle() string
	SetTitle(string)
	GetOriginalURL() string
	SetOriginalURL(string)
	GetMetaMap() map[string]string
	SetMetaMap(map[string]string)
}

type LinkSetActionGetter interface {
	GetAction() shared.ActionType
	GetLinkIds() ([]int64, error)
}
