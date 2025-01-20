package storage

import (
	"context"
	"net/url"

	"github.com/mikeschinkel/savetabs/daemon/shared"
	"github.com/mikeschinkel/savetabs/daemon/sqlc"
)

type UnparsedLink struct {
	OriginalURL string
	Title       string
	URL         *url.URL
}

func LatestUnparsedLinks(ctx context.Context) (links []UnparsedLink, err error) {
	var rows []sqlc.ListLatestUnparsedLinkURLsRow
	err = ExecWithNestedTx(func(dbtx *NestedDBTX) (err error) {
		q := dbtx.DataStore.Queries(dbtx)
		rows, err = q.ListLatestUnparsedLinkURLs(ctx, sqlc.ListLatestUnparsedLinkURLsParams{
			LinksArchived: ArchivedOrNot,
			LinksDeleted:  ArchivedOrNot,
		})
		if err != nil {
			goto end
		}
	end:
		return err
	})
	links = shared.ConvertSlice(rows, func(row sqlc.ListLatestUnparsedLinkURLsRow) UnparsedLink {
		return UnparsedLink{
			OriginalURL: row.OriginalUrl,
		}
	})
	return links, err
}

type ParsedLink struct {
	Title       string
	Scheme      string
	Subdomain   string
	SLD         string
	TLD         string
	Port        string
	Path        string
	Query       string
	Fragment    string
	OriginalURL string
}

func UpdateUnparsedLink(ctx Context, link ParsedLink) error {
	return ExecWithNestedTx(func(dbtx *NestedDBTX) error {
		return dbtx.DataStore.Queries(dbtx).UpdateParsedLink(ctx, sqlc.UpdateParsedLinkParams{
			Title:       link.Title,
			Scheme:      link.Scheme,
			Subdomain:   link.Subdomain,
			Sld:         link.SLD,
			Tld:         link.TLD,
			Port:        link.Port,
			Path:        link.Path,
			Query:       link.Query,
			Fragment:    link.Fragment,
			OriginalUrl: link.OriginalURL,
		})
	})
}
