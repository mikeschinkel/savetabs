package restapi

import (
	"savetabs/storage"
)

var _ storage.LinkWithGroupGetSetter = (*LinkWithGroup)(nil)

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
func (link LinkWithGroup) GetContent() string {
	if link.Html == nil {
		return ""
	}
	return *link.Html
}
func (link LinkWithGroup) SetContent(s string) {
	link.Html = &s
}
func (link LinkWithGroup) GetOriginalURL() string {
	if link.OriginalUrl == nil {
		return ""
	}
	return *link.OriginalUrl
}
func (link LinkWithGroup) SetOriginalURL(s string) {
	link.OriginalUrl = &s
}
func (link LinkWithGroup) GetMetaMap() map[string]string {
	panic("LinkWithGroup.GetMeta() not yet implemented")
	return nil
}
func (link LinkWithGroup) SetMetaMap(m map[string]string) {
	panic("LinkWithGroup.SetMeta() not yet implemented")
}
