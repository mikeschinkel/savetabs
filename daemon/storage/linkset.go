package storage

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mikeschinkel/savetabs/daemon/shared"
	"github.com/mikeschinkel/savetabs/daemon/sqlc"
)

type UpsertLinkset struct {
	Action  shared.Action
	LinkIds []int64
}

func LinksetUpsert(ctx Context, ds DataStore, ls UpsertLinkset) (err error) {

	err = GetNestedDBTX(ds).Exec(func(dbtx *NestedDBTX) (err error) {
		q := ds.Queries(dbtx)
		switch ls.Action {
		case shared.ArchiveAction:
			err = q.ArchiveLinks(ctx, ls.LinkIds)
			if err != nil {
				err = errors.Join(ErrFailedToArchiveLinks, err)
				goto end
			}
		case shared.DeleteAction:
			err = q.MarkLinksDeleted(ctx, ls.LinkIds)
			if err != nil {
				err = errors.Join(ErrFailedToMarkLinksDeleted, err)
				goto end
			}
		}
	end:
		return err
	})
	if err != nil {
		err = errors.Join(err, fmt.Errorf("link_ids=%v", ls.LinkIds))
	}
	return err
}

type LoadLinksetParams struct {
	shared.FilterQuery
}

type LinksetToLoad struct {
	Links  []Link
	Params LoadLinksetParams
}

func LoadLinkset(ctx Context, params LoadLinksetParams) (ls LinksetToLoad, err error) {
	var links []sqlc.ListFilteredLinksRow
	var ids []int64
	var linkIds []int64

	ls.Params = params
	db := GetNestedDBTX(GetDatastore())
	err = db.Exec(func(dbtx *NestedDBTX) (err error) {
		q := dbtx.DataStore.Queries(dbtx)
		me := shared.MultiErr{}
		for _, fi := range params.FilterQuery.FilterItems {
			switch fi.FilterType() {
			case shared.GroupTypeFilterType:
				ids, err = q.ListLinkIdsByGroupType(ctx, sqlc.ListLinkIdsByGroupTypeParams{
					GroupTypes: shared.ConvertSlice(fi.Filters(), func(f any) string {
						return f.(string)
					}),
					LinksArchived: NotArchived,
					LinksDeleted:  NotDeleted,
				})
			case shared.GroupFilterType:
				ids, err = q.ListLinkIdsByGroup(ctx, sqlc.ListLinkIdsByGroupParams{
					Slugs: shared.ConvertSlice(fi.Filters(), func(f any) string {
						return f.(string)
					}),
					LinksArchived: NotArchived,
					LinksDeleted:  NotDeleted,
				})
			case shared.MetaFilterType:
				ids, err = q.ListLinkIdsByMeta(ctx, sqlc.ListLinkIdsByMetaParams{
					Keys: shared.ConvertSlice(fi.Filters(), func(f any) string {
						s := f.(string)
						return s[:strings.Index(s, "=")]
					}),
					KvPairs: shared.ConvertSlice(fi.Filters(), func(f any) string {
						return f.(string)
					}),
					LinksArchived: NotArchived,
					LinksDeleted:  NotDeleted,
				})
			default:
				panicf("Unhandled FilterQuery Type: %v", fi.FilterType())
			}
			if err != nil {
				me.Add(err)
			}
			if len(ids) == 0 {
				continue
			}
			linkIds = append(linkIds, ids...)
		}
		if len(linkIds) == 0 {
			goto end
		}
		links, err = q.ListFilteredLinks(ctx, sqlc.ListFilteredLinksParams{
			Ids:           linkIds,
			LinksArchived: NotArchived,
			LinksDeleted:  NotDeleted,
		})
		if err != nil {
			goto end
		}
		ls.Links = make([]Link, len(links))
		for i, link := range links {
			ls.Links[i] = NewLoadLink(sqlc.LoadLinkRow{
				ID:          link.ID,
				OriginalUrl: link.OriginalUrl,
				Title:       link.Title,
				Scheme:      link.Scheme,
				Subdomain:   link.Subdomain,
				Tld:         link.Tld,
				Sld:         link.Sld,
				Path:        link.Path,
				Query:       link.Query,
				Fragment:    link.Fragment,
				Port:        link.Port,
			})
		}
	end:
		return err
	})
	return ls, err
}

//// TODO: This is for Caretaker task
//func linkFromSetLink(sl sqlc.ListFilteredLinksRow) (link Link) {
//	title := sl.Title
//	u, err := url.Parse(sl.OriginalUrl)
//	if err != nil {
//		title = "ERROR: " + err.Error()
//	}
//	link = NewLoadLink(u)
//	link.Id = sl.ID
//	link.Scheme = title
//	link.Scheme = sl.Scheme
//	link.Subdomain = sl.Subdomain
//	link.SLD = sl.Sld
//	link.TLD = sl.Tld
//	link.Port = sl.Port
//	link.Path = sl.Path
//	link.Query = sl.Query
//	link.Fragment = sl.Fragment
//	return link
//}
