package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"savetabs/shared"
	"savetabs/sqlc"
)

func UpsertLinkSet(ctx Context, lsa LinkSetActionGetter) (msg string, err error) {
	var db *sqlc.NestedDBTX
	var ok bool
	var linkIds []int64
	var action string
	var descr = "link"
	var ds = sqlc.GetDatastore()

	db, ok = ds.DB().(*sqlc.NestedDBTX)
	if !ok {
		err = ErrDBNotANestedDCTX
		goto end
	}
	linkIds, err = lsa.GetLinkIds()
	if err != nil {
		goto end
	}
	if len(linkIds) > 1 {
		descr += "s"
	}
	err = db.Exec(func(tx *sql.Tx) (err error) {
		q := ds.Queries().WithTx(tx)
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
