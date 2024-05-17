package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"

	"savetabs/shared"
)

type LinkWithGroup struct {
	URL         *url.URL
	OriginalURL string
	Title       string
	GroupId     int64
	GroupSlug   string
	GroupType   shared.GroupType
	Group       string
}

type LinkToAdd struct {
	OriginalURL string `json:"original_url"`
	Title       string `json:"title"`
}

type UpsertLinksWithGroupsParams struct {
	Links        []LinkToAdd
	Groups       []Group
	GroupedLinks []LinkGroup
	dbtx         *NestedDBTX
}

func UpsertLinksWithGroups(ctx context.Context, p UpsertLinksWithGroupsParams) error {
	var groupBytes []byte
	var groupedLinkBytes []byte
	var gg []Group
	var rgs []LinkGroup
	//var mm []Meta
	var me = shared.NewMultiErr()

	slog.Info("Received from Chrome extension", "num_links", len(p.Links))

	linkBytes, err := json.Marshal(p.Links)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "link"))
	}

	groupBytes, err = json.Marshal(p.Groups)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "group"))
	}

	groupedLinkBytes, err = json.Marshal(p.GroupedLinks)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "link_group"))
	}

	err = execWithEnsuredNestedDBTX(p.dbtx, func(dbtx *NestedDBTX) error {
		var innerME = shared.NewMultiErr()

		err := UpsertLinksFromJSON(ctx, dbtx, string(linkBytes))
		if err != nil {
			innerME.Add(err, ErrFailedUpsertLinks)
		}

		err = UpsertLinkGroupsFromJSON(ctx, dbtx, string(groupedLinkBytes))
		if err != nil {
			innerME.Add(err, ErrFailedUpsertLinkGroups)
		}

		err = UpsertGroupsFromJSON(ctx, dbtx, string(groupBytes))
		if err != nil {
			innerME.Add(err, ErrFailedUpsertGroups)
		}

		slog.Info("Saved",
			"links", len(p.Links),
			"link_groups", len(rgs),
			"groups", len(gg),
		)
		throttle()
		return innerME.Err()
	})
	return me.Err()
}
