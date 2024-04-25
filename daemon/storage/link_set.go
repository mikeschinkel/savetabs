package storage

import (
	"database/sql"
	"fmt"

	"savetabs/sqlc"
)

func UpsertLinkSet(ctx Context, lsa LinkSetActionGetter) (err error) {
	var ds sqlc.DataStore
	var db *sqlc.NestedDBTX
	var ok bool

	ds = sqlc.GetDatastore()
	db, ok = ds.DB().(*sqlc.NestedDBTX)
	if !ok {
		err = ErrDBNotANestedDCTX
		goto end
	}
	err = db.Exec(func(tx *sql.Tx) (err error) {
		fmt.Printf("Action: %s\n", lsa.GetAction())
		linkIds, err := lsa.GetLinkIds()
		if err != nil {
			goto end
		}
		fmt.Printf("Links: %v\n", linkIds)
	end:
		return err
	})
end:
	return err
}
