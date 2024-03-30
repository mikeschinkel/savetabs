package sqlc

import (
	"context"
	"database/sql"
)

func UpsertResources(ctx context.Context, ds DataStore, resourcesJSON string) (err error) {
	return upsertFromJSON(ctx, ds, resourcesJSON, func(ctx context.Context, q *Queries, id int64) error {
		return q.UpsertResourcesFromVarJSON(ctx, id)
	})
}

func UpsertGroups(ctx context.Context, ds DataStore, groupsJSON string) (err error) {
	return upsertFromJSON(ctx, ds, groupsJSON, func(ctx context.Context, q *Queries, id int64) error {
		return q.UpsertGroupsFromVarJSON(ctx, id)
	})
}

func UpsertKeyValues(ctx context.Context, ds DataStore, keyValuesJSON string) (err error) {
	return upsertFromJSON(ctx, ds, keyValuesJSON, func(ctx context.Context, q *Queries, id int64) error {
		return q.UpsertKeyValuesFromVarJSON(ctx, id)
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

		q := ds.Queries()
		varId, err = q.WithTx(tx).UpsertVar(ctx, UpsertVarParams{
			Key:   "json",
			Value: NewNullString(j),
		})
		if err != nil {
			goto end
		}
		err = fn(ctx, q.WithTx(tx), varId)
		if err != nil {
			goto end
		}
		err = q.WithTx(tx).DeleteVar(ctx, varId)
		if err != nil {
			goto end
		}
	end:
		return err
	})
end:
	return err
}
