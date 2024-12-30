package storage

import (
	"context"
)

var ds DataStore

func GetDatastore() DataStore {
	return ds
}

type Args struct {
	AppName string
}

func Initialize(ctx context.Context, args Args) (err error) {

	if len(args.AppName) == 0 {
		err = ErrAppNameMustNotBeEmpty
		goto end
	}
	ds = NewSqliteDataStore(args)

end:
	return ds.Initialize(ctx)
}
