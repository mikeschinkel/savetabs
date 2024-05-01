package storage

import (
	"errors"
	"fmt"

	"savetabs/shared"
	"savetabs/sqlc"
)

func UpsertLinkSet(ctx Context, ds sqlc.DataStore, lsa LinkSetActionGetter) (msg string, err error) {
	var db *sqlc.NestedDBTX
	var linkIds []int64
	var action string
	var descr = "link"

	db = sqlc.GetNestedDBTX(ds)
	linkIds, err = lsa.GetLinkIds()
	if err != nil {
		goto end
	}
	if len(linkIds) > 1 {
		descr += "s"
	}
	err = db.Exec(func(tx sqlc.DBTX) (err error) {
		q := ds.Queries(tx)
		switch lsa.GetAction() {
		case shared.ArchiveAction:
			err = execForLinkdIds(ctx, q, lsa, func(context Context, int64s []int64) error {
				return q.ArchiveLinks(ctx, linkIds)
			})
			if err != nil {
				goto end
			}
			action = "archived"
		case shared.DeleteAction:
			err = execForLinkdIds(ctx, q, lsa, func(context Context, int64s []int64) error {
				return q.DeleteLinks(ctx, linkIds)
			})
			if err != nil {
				goto end
			}
			action = "deleted"
		}
	end:
		return err
	})
	msg = fmt.Sprintf("%d %s %s", len(linkIds), descr, action)
end:
	return msg, err
}

func execForLinkdIds(ctx Context, q *sqlc.Queries, lsa LinkSetActionGetter, execFn func(Context, []int64) error) (err error) {
	var linkIds []int64
	linkIds, err = lsa.GetLinkIds()
	if err != nil {
		err = errors.Join(err, fmt.Errorf("link_ids=%v", linkIds))
		goto end
	}
	err = execFn(ctx, linkIds)
	if err != nil {
		err = errors.Join(ErrFailedToArchiveLinks, fmt.Errorf("link_ids=%v", linkIds))
		goto end
	}
end:
	return err
}
