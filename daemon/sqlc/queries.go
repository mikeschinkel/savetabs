package sqlc

import (
	"context"
)

func UpsertLinks(ctx context.Context, db *NestedDBTX, linksJSON string) (err error) {
	return upsertFromJSON(ctx, db, linksJSON, func(ctx context.Context, q *Queries, id int64) error {
		return q.UpsertLinksFromVarJSON(ctx, id)
	})
}

func UpsertLinkMeta(ctx context.Context, db *NestedDBTX, metaJSON string) (err error) {
	return upsertFromJSON(ctx, db, metaJSON, func(ctx context.Context, q *Queries, id int64) error {
		return q.UpsertLinkMetaFromVarJSON(ctx, id)
	})
}

func UpsertGroups(ctx context.Context, db *NestedDBTX, groupsJSON string) (err error) {
	return upsertFromJSON(ctx, db, groupsJSON, func(ctx context.Context, q *Queries, id int64) error {
		return q.UpsertGroupsFromVarJSON(ctx, id)
	})
}

func UpsertLinkGroups(ctx context.Context, db *NestedDBTX, rgsJSON string) (err error) {
	return upsertFromJSON(ctx, db, rgsJSON, func(ctx context.Context, q *Queries, id int64) error {
		return q.UpsertLinkGroupsFromVarJSON(ctx, id)
	})
}

func UpsertMeta(ctx context.Context, db *NestedDBTX, metaJSON string) (err error) {
	return upsertFromJSON(ctx, db, metaJSON, func(ctx context.Context, q *Queries, id int64) error {
		return q.UpsertMetaFromVarJSON(ctx, id)
	})
}

func upsertFromJSON(ctx context.Context, db *NestedDBTX, j string, fn func(ctx context.Context, q *Queries, id int64) error) (err error) {
	err = db.Exec(func(dbtx DBTX) (err error) {
		var varId int64

		q := ds.Queries(dbtx)
		varId, err = q.UpsertVar(ctx, UpsertVarParams{
			Key:   "json",
			Value: NewNullString(j),
		})
		if err != nil {
			goto end
		}
		err = fn(ctx, q, varId)
		if err != nil {
			goto end
		}
		err = q.DeleteVar(ctx, varId)
		if err != nil {
			goto end
		}
	end:
		return err
	})

	return err
}
