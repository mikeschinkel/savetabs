package model

import (
	"net/url"

	"savetabs/storage"
)

type LinkToAdd struct {
	Id          int64
	URL         url.URL
	OriginalURL string
	Title       string
	Content     Content
}

func (ul LinkToAdd) Body() string {
	return ul.Content.Body.String()
}

func (ul LinkToAdd) Head() string {
	return ul.Content.Head.String()
}

func AddLink(ctx Context, link LinkToAdd) (linkId int64, err error) {
	err = storage.ExecWithNestedTx(func(dbtx *storage.NestedDBTX) (err error) {
		linkId, err = storage.LinkUpsert(ctx, dbtx, storage.UpsertLink{
			Id:          link.Id,
			URL:         link.URL,
			OriginalURL: link.OriginalURL,
			Title:       link.Title,
			Head:        link.Head(),
			Body:        link.Body(),
		})
		return err
	})
	return linkId, err
}

type Link storage.Link

func LoadLink(ctx Context, linkId int64) (link Link, err error) {
	err = storage.ExecWithNestedTx(func(dbtx *storage.NestedDBTX) error {
		var ll storage.Link
		ll, err = storage.LinkLoad(ctx, linkId)
		if err != nil {
			goto end
		}
		link = Link(ll)
	end:
		return err
	})
	return link, err
}
