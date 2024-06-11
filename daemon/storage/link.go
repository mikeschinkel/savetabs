package storage

import (
	"errors"
	"fmt"
	"net/url"
	"slices"

	"savetabs/sqlc"
)

type UpsertLink struct {
	Id          int64
	URL         url.URL
	OriginalURL string
	Title       string
	Head        string
	Body        string
}

func LinkUpsert(ctx Context, db *NestedDBTX, link UpsertLink) (linkId int64, err error) {
	err = db.Exec(func(dbtx *NestedDBTX) (err error) {
		q := db.DataStore.Queries(dbtx)
		linkId, err = q.UpsertLink(ctx, sqlc.UpsertLinkParams{
			OriginalUrl: link.OriginalURL,
			Title:       link.Title,
		})
		if err != nil {
			err = errors.Join(ErrFailedUpsertLink, err)
			goto end
		}
		if linkId == 0 {
			linkId, err = q.LoadLinkIdByUrl(ctx, link.OriginalURL)
		}
		if err != nil {
			err = errors.Join(
				ErrFailedLoadLinkByUrl,
				fmt.Errorf("url=%s", link.OriginalURL),
				err,
			)
			goto end
		}
		err = InsertContent(ctx, dbtx, ContentToInsert{
			LinkId: linkId,
			Head:   link.Head,
			Body:   link.Body,
		})
		if err != nil {
			goto end
		}
		//mj = link.MetaJSON()
		//err = sqlc.UpsertLinkMetaFromJSON(ctx, db, mj)
		//if err != nil {
		//	err = errors.Join(ErrFailedUpsertLinkMeta, fmt.Errorf("meta_json=%s", mj), err)
		//	goto end
		//}
	end:
		return err
	})
	return linkId, err
}

type Link struct {
	Id        int64             `json:"id"`
	URL       string            `json:"url"`
	Title     string            `json:"title"`
	Scheme    string            `json:"scheme"`
	Host      string            `json:"host"`
	Subdomain string            `json:"subdomain"`
	TLD       string            `json:"tld"`
	SLD       string            `json:"sld"`
	Port      string            `json:"port"`
	Path      string            `json:"path"`
	Query     string            `json:"query"`
	Fragment  string            `json:"fragment"`
	MetaMap   map[string]string `json:"meta"`
}

func LinkLoad(ctx Context, linkId int64) (link Link, err error) {
	db := GetNestedDBTX(GetDatastore())
	err = db.Exec(func(dbtx *NestedDBTX) error {
		var q = db.DataStore.Queries(dbtx)
		row, err := q.LoadLink(ctx, linkId)
		if err != nil {
			goto end
		}
		link = NewLoadLink(row)
	end:
		return err
	})
	return link, err
}

func NewLoadLink(row sqlc.LoadLinkRow) Link {
	return Link{
		Id:        row.ID,
		URL:       row.OriginalUrl,
		Title:     row.Title,
		Scheme:    row.Scheme,
		Subdomain: row.Subdomain,
		TLD:       row.Tld,
		SLD:       row.Sld,
		Port:      row.Port,
		Path:      row.Path,
		Query:     row.Query,
		Fragment:  row.Fragment,
	}
}

func ValidateLinks(ctx Context, dbtx *NestedDBTX, linkIds []int64) (missing []int64, err error) {
	err = execWithEnsuredNestedDBTX(dbtx, func(dbtx *NestedDBTX) error {

		q := dbtx.DataStore.Queries(dbtx)
		ids, err := q.ValidateLinks(ctx, linkIds)
		missing = make([]int64, 0, len(ids))
		for _, id := range linkIds {
			if slices.Contains(ids, id) {
				continue
			}
			missing = append(missing, id)
		}
		return err
	})
	return missing, err
}
