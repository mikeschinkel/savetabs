package restapi

import (
	"savetabs/storage"
)

var _ storage.LinkWithGroupPropGetSetter = (*LinkWithGroup)(nil)

func (link LinkWithGroup) GetGroup() string {
	if link.Group == nil {
		return ""
	}
	return *link.Group
}
func (link LinkWithGroup) SetGroup(s string) {
	link.Group = &s
}
func (link LinkWithGroup) GetGroupId() int64 {
	if link.GroupId == nil {
		return 0
	}
	return *link.GroupId
}
func (link LinkWithGroup) SetGroupId(id int64) {
	link.GroupId = &id
}
func (link LinkWithGroup) GetGroupType() string {
	if link.GroupType == nil {
		return ""
	}
	return *link.GroupType
}
func (link LinkWithGroup) SetGroupType(s string) {
	link.GroupType = &s
}
func (link LinkWithGroup) GetId() int64 {
	if link.Id == nil {
		return 0
	}
	return *link.Id
}
func (link LinkWithGroup) SetId(id int64) {
	link.Id = &id
}
func (link LinkWithGroup) GetTitle() string {
	if link.Title == nil {
		return ""
	}
	return *link.Title
}
func (link LinkWithGroup) SetTitle(s string) {
	link.Title = &s
}
func (link LinkWithGroup) GetURL() string {
	if link.Url == nil {
		return ""
	}
	return *link.Url
}
func (link LinkWithGroup) SetURL(s string) {
	link.Url = &s
}
