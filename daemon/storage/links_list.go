package storage

import (
	"savetabs/shared"
	"savetabs/sqlc"
)

type ListLinksArgs sqlc.ListLinksParams

func ListLinks(ctx Context, p ListLinksArgs) ([]Link, error) {
	sqlcLinks := make([]sqlc.Link, 0)
	err := ExecWithNestedTx(func(dbtx *NestedDBTX) (err error) {
		q := dbtx.DataStore.Queries(dbtx)
		sqlcLinks, err = q.ListLinks(ctx, sqlc.ListLinksParams(p))
		return err
	})
	return shared.ConvertSlice(sqlcLinks, func(link sqlc.Link) Link {
		return Link{
			Id:        link.ID,
			URL:       link.OriginalUrl,
			Title:     link.Title,
			Scheme:    link.Scheme,
			Host:      link.Host,
			Subdomain: link.Subdomain,
			TLD:       link.Tld,
			SLD:       link.Sld,
			Port:      link.Port,
			Path:      link.Path,
			Query:     link.Query,
			Fragment:  link.Fragment,
		}
	}), err
}
