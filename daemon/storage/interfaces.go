package storage

type StorageModifier interface {
	UpsertLinksWithGroups(Context, LinksWithGroups) error
}

type LinksWithGroupsGetSetter interface {
	GetLinkCount() int
	GetLinksWithGroups() []LinkWithGroupPropGetSetter
	SetLinksWithGroups([]LinkWithGroupPropGetSetter)
}

type LinkWithGroupPropGetSetter interface {
	GetGroup() string
	SetGroup(string)
	GetGroupId() int64
	SetGroupId(int64)
	GetGroupType() string
	SetGroupType(string)
	GetId() int64
	SetId(int64)
	GetTitle() string
	SetTitle(string)
	GetURL() string
	SetURL(string)
}
