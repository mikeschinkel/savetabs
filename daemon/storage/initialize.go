package storage

import (
	"context"
)

var ds DataStore

func GetDatastore() DataStore {
	return ds
}

func Initialize(ctx context.Context, fp string) (err error) {

	ds = NewSqliteDataStore(fp)

	err = ds.Initialize(ctx)
	if err != nil {
		goto end
	}
end:
	if err != nil {
		err = Error(ErrFailedToInitDataStore, err, "data_file", fp)
	}
	return err
}
