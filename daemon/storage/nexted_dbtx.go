package storage

import (
	"database/sql"
	"errors"
	"log/slog"

	"savetabs/sqlc"
)

type NestedDBTX struct {
	sqlc.DBTX
	tx        *sql.Tx
	level     int
	DataStore DataStore
}

func NewNestedDBTX(ds DataStore) *NestedDBTX {
	return &NestedDBTX{
		DataStore: ds,
		DBTX:      ds.DB(),
	}
}

func (dbtx *NestedDBTX) Begin() (_ *sql.Tx, err error) {
	var db *sql.DB
	var ok bool
	dbtx.level++
	if dbtx.level > 1 {
		goto end
	}
	db, ok = dbtx.DBTX.(*sql.DB)
	if !ok {
		panic("dbtx not a *sql.DB")
	}
	dbtx.tx, err = db.Begin()
	if err != nil {
		goto end
	}
end:
	return dbtx.tx, err
}

func (dbtx *NestedDBTX) Commit() (err error) {
	dbtx.level--
	if dbtx.level > 0 {
		return nil
	}
	err = dbtx.tx.Commit()
	dbtx.tx = nil
	return err
}

func (dbtx *NestedDBTX) Rollback() (err error) {
	dbtx.level--
	if dbtx.level > 0 {
		return nil
	}
	err = dbtx.tx.Rollback()
	dbtx.tx = nil
	return err
}

// Exec is a method which abstracts a sql transaction
//
// See: https://stackoverflow.com/a/44522218/102699
func (dbtx *NestedDBTX) Exec(fn func(dbtx *NestedDBTX) error) (err error) {
	_, err = dbtx.Begin()
	defer func() {
		if err != nil {
			slog.Error("NestedDBTX", "error", err.Error())
			_err := dbtx.Rollback()
			if _err != nil {
				// If this happens, figure out how to make more robust
				slog.Error(_err.Error())
				err = errors.Join(err, _err)
			}
			return
		}
		err = dbtx.Commit()
	}()
	if err != nil {
		return err
	}

	// run transaction
	err = fn(dbtx)

	return err
}

func GetNestedDBTX(ds DataStore) *NestedDBTX {
	return NewNestedDBTX(ds)
}
