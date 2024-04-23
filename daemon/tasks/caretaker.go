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
	Queries   *sqlc.Queries
}

func NewCaretaker(ds sqlc.DataStore) *Caretaker {
	return &Caretaker{
		DataStore: ds,
		Queries:   ds.Queries(),
	}
}

func (c Caretaker) Run(ctx context.Context) (err error) {
	var ll []sqlc.ListLatestUnparsedLinkURLsRow
	var u *url.URL

	slog.Info("Running Caretaker")
	ll, err = c.Queries.ListLatestUnparsedLinkURLs(ctx, sqlc.NotArchived)
	if err != nil {
		goto end
	}
	for _, link := range ll {
		slog.Info("Processing", "url", link.OriginalUrl) // TODO: Change to slog.Debug()
		u, err = url.Parse(link.OriginalUrl)
		if err != nil {
			slog.Error(err.Error(), "url", link.OriginalUrl)
			continue
		}
		var host Host
		host = parseHost(u)

		parts := sqlc.UpdateLinkPartsParams{
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
		err = c.Queries.UpdateLinkParts(ctx, parts)
		if err != nil {
			slog.Error(err.Error(), "url", link.OriginalUrl)
			continue
		}
		// Allow for some CPU breathing room
		time.Sleep(3 * time.Second)
	}
end:
	return err
}
