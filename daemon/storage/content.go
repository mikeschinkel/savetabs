package storage

import (
	"errors"

	"savetabs/sqlc"
)

type ContentToInsert struct {
	Id     int64
	LinkId int64
	Head   string
	Body   string
}

func InsertContent(ctx Context, dbtx *NestedDBTX, content ContentToInsert) (err error) {
	err = execWithEnsuredNestedDBTX(dbtx, func(dbtx *NestedDBTX) error {
		q := dbtx.DataStore.Queries(dbtx)
		return q.InsertContent(ctx, sqlc.InsertContentParams{
			LinkID: content.LinkId,
			Head:   content.Head,
			Body:   content.Body,
		})
	})
	if err != nil {
		err = errors.Join(ErrFailedInsertLinkContent, err)
		goto end
	}
end:
	return err
}
