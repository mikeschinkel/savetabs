package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"savetabs/sqlc"
	"savetabs/storage"
)

var _ storage.LinksWithGroupsGetter = (*linksWithGroups)(nil)

type linksWithGroupsForJSON []linkWithGroupForJSON

func (links linksWithGroupsForJSON) GetLinkCount() int {
	return len(links)
}

func (links linksWithGroupsForJSON) GetLinksWithGroups() storage.LinksWithGroupsGetter {
	ll := make(linksWithGroups, links.GetLinkCount())
	for i, link := range links {
		ll[i] = LinkWithGroup{
			Group:       &link.Group,
			GroupId:     &link.GroupId,
			GroupType:   &link.GroupType,
			OriginalUrl: &link.OriginalURL,
			Title:       &link.Title,
		}
	}
	return ll
}

type linkWithGroupForJSON struct {
	OriginalURL string `json:"original_url"`
	Title       string `json:"title"`
	GroupId     int64  `json:"groupId"`
	GroupType   string `json:"groupType"`
	Group       string `json:"group"`
}

func (a *API) PostLinksWithGroups(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Find a better status result than "Bad Gateway"
		a.sendError(w, r, http.StatusBadGateway, err.Error())
		return
	}
	var links linksWithGroupsForJSON
	err = json.Unmarshal(body, &links)
	if err != nil {
		a.sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	db := sqlc.GetNestedDBTX(sqlc.GetDatastore())
	err = db.Exec(func(dbtx sqlc.DBTX) error {
		return storage.UpsertLinksWithGroups(ctx, db, links.GetLinksWithGroups())
	})
	switch {
	case err == nil:
		goto end
	case errors.Is(err, ErrFailedToUnmarshal):
		a.sendError(w, r, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrFailedUpsertLinks):
		// TODO: Break out errors into different status for different reasons
		fallthrough
	default:
		a.sendError(w, r, http.StatusInternalServerError, err.Error())
	}
end:
}

var _ storage.LinksWithGroupsGetter = (*linksWithGroups)(nil)

type linksWithGroups LinksWithGroups

func (ls linksWithGroups) GetLinkCount() int {
	return len(ls)
}

func (ls linksWithGroups) GetLinksWithGroups() []storage.LinkWithGroupGetSetter {
	ll := make([]storage.LinkWithGroupGetSetter, ls.GetLinkCount())
	for i, link := range ls {
		ll[i] = link
	}
	return ll
}
