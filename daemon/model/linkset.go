package model

import (
	"github.com/mikeschinkel/savetabs/daemon/storage"
)

type LinksetToAdd storage.UpsertLinkset

func AddLinksetIfNotExist(ctx Context, ls LinksetToAdd) (err error) {
	return storage.ExecWithNestedTx(func(dbtx *storage.NestedDBTX) error {
		return storage.LinksetUpsert(ctx, dbtx.DataStore, storage.UpsertLinkset(ls))
	})
}

func LoadLinkURLs(ctx Context, linkIds []int64) (urls []string, err error) {
	ds := storage.GetDatastore()
	db := storage.GetNestedDBTX(ds)
	// TODO: Load from model, not storage
	return ds.Queries(db).GetLinkURLs(ctx, linkIds)
}

type LinksetToLoadParams storage.LoadLinksetParams

type LinksetToLoad struct {
	Links  []Link
	Params LinksetToLoadParams
}

func LoadLinkset(ctx Context, params LinksetToLoadParams) (_ LinksetToLoad, err error) {
	var ls storage.LinksetToLoad
	ls, err = storage.LoadLinkset(ctx, storage.LoadLinksetParams(params))
	if err != nil {
		goto end
	}
end:
	return NewLinksetToLoad(ls), err
}

func NewLinksetToLoad(ls storage.LinksetToLoad) LinksetToLoad {
	links := make([]Link, len(ls.Links))
	for i, l := range ls.Links {
		links[i] = Link(l)
	}
	return LinksetToLoad{
		Links:  links,
		Params: LinksetToLoadParams(ls.Params),
	}
}
