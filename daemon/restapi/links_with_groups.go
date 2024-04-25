package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"savetabs/storage"
)

func (a *API) PostLinksWithGroups(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Find a better status result than "Bad Gateway"
		a.sendError(w, r, http.StatusBadGateway, err.Error())
		return
	}
	var links LinksWithGroups
	err = json.Unmarshal(body, &links)
	if err != nil {
		a.sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	err = storage.UpsertLinksWithGroups(ctx, linksWithGroups(links))
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

var _ storage.LinksWithGroupsGetSetter = (*linksWithGroups)(nil)

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
