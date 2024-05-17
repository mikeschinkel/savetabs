package storage

func ExecWithNestedTx(fn func(dbtx *NestedDBTX) (err error)) error {
	db := GetNestedDBTX(GetDatastore())
	return db.Exec(fn)
}
