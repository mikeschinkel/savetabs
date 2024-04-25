package sqlc

import (
	"context"
	"database/sql"
)

func UpsertLinks(ctx context.Context, ds DataStore, linksJSON string) (err error) {
	return upsertFromJSON(ctx, ds, linksJSON, func(ctx context.Context, q *Queries, id int64) error {
		return q.UpsertLinksFromVarJSON(ctx, id)
	})
}

func UpsertGroups(ctx context.Context, ds DataStore, groupsJSON string) (err error) {
	return upsertFromJSON(ctx, ds, groupsJSON, func(ctx context.Context, q *Queries, id int64) error {
		return q.UpsertGroupsFromVarJSON(ctx, id)
	})
}

func UpsertLinkGroups(ctx context.Context, ds DataStore, rgsJSON string) (err error) {
	return upsertFromJSON(ctx, ds, rgsJSON, func(ctx context.Context, q *Queries, id int64) error {
		return q.UpsertLinkGroupsFromVarJSON(ctx, id)
	})
}

func UpsertMetadata(ctx context.Context, ds DataStore, metadataJSON string) (err error) {
	return upsertFromJSON(ctx, ds, metadataJSON, func(ctx context.Context, q *Queries, id int64) error {
		return q.UpsertMetadataFromVarJSON(ctx, id)
	})
}

func upsertFromJSON(ctx context.Context, ds DataStore, j string, fn func(ctx context.Context, q *Queries, id int64) error) (err error) {
	db, ok := ds.DB().(*NestedDBTX)
	if !ok {
		err = ErrDBNotANestedDBTX
		goto end
	}
	err = db.Exec(func(tx *sql.Tx) (err error) {
		var varId int64

		q := ds.Queries().WithTx(tx)
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
end:
	return err
}
