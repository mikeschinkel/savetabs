package sqlc

import (
	"database/sql"
	"log"
)

type NestedDBTX struct {
	DBTX
	tx    *sql.Tx
	level int
}

func NewNestedDBTX(db DBTX) DBTX {
	return &NestedDBTX{
		DBTX: db,
	}
}

func (dbtx *NestedDBTX) Begin() (tx *sql.Tx, err error) {
	var db *sql.DB
	var ok bool

	if dbtx.level > 0 {
		dbtx.level++
		goto end
	}
	db, ok = dbtx.DBTX.(*sql.DB)
	if !ok {
		panic("dbtx not a *sql.DB")
	}
	dbtx.tx, err = db.Begin()
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
func (dbtx *NestedDBTX) Exec(fn func(*sql.Tx) error) (err error) {
	var tx *sql.Tx
	tx, err = dbtx.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			err = tx.Rollback()
			if err != nil {
				// If this happens, figure out how to make more robust
				log.Fatal(err.Error())
			}
			return
		}
		err = tx.Commit()
	}()

	// run transaction
	err = fn(tx)

	return err
}
