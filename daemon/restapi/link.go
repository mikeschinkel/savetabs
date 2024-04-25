package restapi

import (
	"savetabs/storage"
)

var _ storage.LinkGetSetter = (*Link)(nil)

func (l Link) GetId() int64 {
	return *l.Id
}

func (l Link) SetId(id int64) {
	l.Id = &id
}

func (l Link) GetOriginalURL() string {
	return *l.Url
}

func (l Link) SetOriginalURL(u string) {
	l.Url = &u
}
