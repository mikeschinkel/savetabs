package tasks

import (
	"context"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"savetabs/sqlc"
)

var _ Runner = (*Caretaker)(nil)

type Caretaker struct {
	DataStore sqlc.DataStore
}

func NewCaretaker(ds sqlc.DataStore) *Caretaker {
	return &Caretaker{
		DataStore: ds,
	}
}

func (c Caretaker) Queries(dbtx sqlc.DBTX) *sqlc.Queries {
	return c.DataStore.Queries(dbtx)
}

func (c Caretaker) Run(ctx context.Context) (err error) {
	var ll []sqlc.ListLatestUnparsedLinkURLsRow

	slog.Info("Running Caretaker")
	defer slog.Info("Caretaker run complete")

	db := sqlc.GetNestedDBTX(c.DataStore)
	ll, err = c.Queries(db).ListLatestUnparsedLinkURLs(ctx, sqlc.ListLatestUnparsedLinkURLsParams{})
	if err != nil {
		goto end
	}
	for _, link := range ll {
		err = c.processUnparsedLink(ctx, db, link)
		if err != nil {
			slog.Error("Failed to parse link", "link", link)
		}
		// Allow for some CPU breathing room
		time.Sleep(3 * time.Second)
	}
end:
	return err
}

func (c Caretaker) processUnparsedLink(ctx Context, db *sqlc.NestedDBTX, link sqlc.ListLatestUnparsedLinkURLsRow) error {
	var u *url.URL
	var host Host
	var parts sqlc.UpdateLinkPartsParams

	return db.Exec(func(dbtx sqlc.DBTX) (err error) {
		slog.Info("Processing", "url", link.OriginalUrl) // TODO: Change to slog.Debug()
		u, err = url.Parse(link.OriginalUrl)
		if err != nil {
			slog.Error(err.Error(), "url", link.OriginalUrl)
			goto end
		}

		host = parseHost(u)
		parts = sqlc.UpdateLinkPartsParams{
			Title:       u.String(), // TODO: Change this to real title
			Scheme:      u.Scheme,
			Subdomain:   host.Subdomain(),
			Sld:         host.Sld,
			Tld:         host.TLD(),
			Port:        u.Port(),
			Path:        u.Path,
			Query:       strings.ReplaceAll(u.RawQuery, "%20", "+"),
			Fragment:    strings.ReplaceAll(u.Fragment, "%20", "+"),
			OriginalUrl: link.OriginalUrl,
		}
		slog.Debug("Updating link", "url", link.OriginalUrl, "parts", parts)
		err = c.Queries(db).UpdateLinkParts(ctx, parts)
		if err != nil {
			slog.Error(err.Error(), "url", link.OriginalUrl)
			goto end
		}
	end:
		return err
	})
}
