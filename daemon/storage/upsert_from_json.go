package storage

import (
	"context"

	"savetabs/sqlc"
)

func UpsertLinksFromJSON(ctx context.Context, db *NestedDBTX, linksJSON string) (err error) {
	return upsertFromJSON(ctx, db, linksJSON, func(ctx context.Context, q *sqlc.Queries, id int64) error {
		return q.UpsertLinksFromVarJSON(ctx, id)
	})
}

func UpsertLinkMetaFromJSON(ctx context.Context, db *NestedDBTX, metaJSON string) (err error) {
	return upsertFromJSON(ctx, db, metaJSON, func(ctx context.Context, q *sqlc.Queries, id int64) error {
		return q.UpsertLinkMetaFromVarJSON(ctx, id)
	})
}

func UpsertGroupsFromJSON(ctx context.Context, db *NestedDBTX, groupsJSON string) (err error) {
	return upsertFromJSON(ctx, db, groupsJSON, func(ctx context.Context, q *sqlc.Queries, id int64) error {
		return q.UpsertGroupsFromVarJSON(ctx, id)
	})
}

func UpsertLinkGroupsFromJSON(ctx context.Context, db *NestedDBTX, rgsJSON string) (err error) {
	return upsertFromJSON(ctx, db, rgsJSON, func(ctx context.Context, q *sqlc.Queries, id int64) error {
		return q.UpsertLinkGroupsFromVarJSON(ctx, id)
	})
}

func UpsertMetaFromJSON(ctx context.Context, db *NestedDBTX, metaJSON string) (err error) {
	return upsertFromJSON(ctx, db, metaJSON, func(ctx context.Context, q *sqlc.Queries, id int64) error {
		return q.UpsertMetaFromVarJSON(ctx, id)
	})
}

func upsertFromJSON(ctx context.Context, db *NestedDBTX, j string, fn func(ctx context.Context, q *sqlc.Queries, id int64) error) (err error) {
	err = db.Exec(func(dbtx *NestedDBTX) (err error) {
		var varId int64
		q := db.DataStore.Queries(dbtx)
		varId, err = q.UpsertVar(ctx, sqlc.UpsertVarParams{
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
