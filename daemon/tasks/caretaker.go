package tasks

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/mikeschinkel/savetabs/daemon/guard"
)

type unparsedLink guard.UnparsedLink

type LinkToParse struct {
	url          *url.URL
	Subdomain    string
	SLD          string
	Port         string
	IsIP         bool
	IsLocal      bool
	HasSubdomain bool
}

func run(ctx context.Context) (err error) {
	var links []guard.UnparsedLink

	slog.Info("Running Caretaker")
	defer slog.Info("Caretaker run complete")

	links, err = guard.LatestUnparsedLinks(ctx)
	if err != nil {
		goto end
	}
	for _, link := range links {
		err = updateUnparsedLink(ctx, unparsedLink(link))
		if err != nil {
			slog.Error("Failed to parse link", "link", link)
		}
		// Allow for some CPU breathing room
		time.Sleep(3 * time.Second)
	}
end:
	return err
}

func updateUnparsedLink(ctx Context, link unparsedLink) (err error) {
	slog.Info("Processing", "url", link.OriginalURL) // TODO: Change to slog.Debug()
	link.URL, err = url.Parse(link.OriginalURL)
	if err != nil {
		slog.Error(err.Error(), "url", link.URL)
		goto end
	}
	slog.Debug("Updating link", "link", link)
	err = guard.UpdateUnparsedLink(ctx, guard.UnparsedLink(link))
	if err != nil {
		slog.Error("Failed to update parsed link.",
			"url", link.OriginalURL,
			"error", err,
		)
		goto end
	}
end:
	return err
}
