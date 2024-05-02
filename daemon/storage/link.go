package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"savetabs/sqlc"
)

type Link struct {
	LinkGetSetter
}

func UpsertLink(ctx Context, db *sqlc.NestedDBTX, lnk LinkGetSetter) (linkId int64, err error) {
	var _link LinkGetSetter

	_link, err = sanitizeLink(lnk)
	if err != nil {
		goto end
	}
	linkId, err = Link{_link}.Upsert(ctx, db)
end:
	return linkId, err
}

func (link Link) Upsert(ctx Context, db *sqlc.NestedDBTX) (linkId int64, err error) {
	var c *Content

	u := link.GetOriginalURL()
	err = db.Exec(func(dbtx sqlc.DBTX) (err error) {
		//var mj string

		q := db.DataStore.Queries(dbtx)
		linkId, err = q.UpsertLink(ctx, sqlc.UpsertLinkParams{
			OriginalUrl: u,
			Title:       link.GetTitle(),
		})
		if err != nil {
			err = errors.Join(ErrFailedUpsertLink, err)
			goto end
		}
		if linkId == 0 {
			linkId, err = q.LoadLinkIdByUrl(ctx, u)
		}
		if err != nil {
			err = errors.Join(ErrFailedLoadLinkByUrl, fmt.Errorf("url=%s", u), err)
			goto end
		}
		link.SetId(linkId)
		c = link.Content()
		err = q.InsertContent(ctx, sqlc.InsertContentParams{
			LinkID: link.GetId(),
			Head:   c.Head,
			Body:   c.Body,
		})
		if err != nil {
			err = errors.Join(ErrFailedUpsertLink, err)
			goto end
		}
		//mj = link.MetaJSON()
		//err = sqlc.UpsertLinkMeta(ctx, db, mj)
		//if err != nil {
		//	err = errors.Join(ErrFailedUpsertLinkMeta, fmt.Errorf("meta_json=%s", mj), err)
		//	goto end
		//}
	end:
		return err
	})
	return linkId, err
}

// Meta returns a slice of Meta from the MetaMap property
func (l Link) Meta() []Meta {
	mm := l.GetMetaMap()
	kvs := make([]Meta, len(mm))
	i := 0
	for k, v := range mm {
		kvs[i] = Meta{
			Key:    k,
			Value:  v,
			LinkId: l.GetId(),
		}
		i++
	}
	return kvs
}

func (l Link) MetaJSON() string {
	mm := l.Meta()
	j, err := json.Marshal(mm)
	if err != nil {
		slog.Error("failed to marshal link meta",
			"error", err,
			"meta_map", l.GetMetaMap(),
			"meta", mm,
		)
	}
	return string(j)
}

func (l Link) Content() *Content {
	c := &Content{
		LinkId: l.GetId(),
	}
	c.setRawContent(l.GetContent())
	return c
}
