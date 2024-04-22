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
		sendError(w, r, http.StatusBadGateway, err.Error())
		return
	}
	var links LinksWithGroups
	err = json.Unmarshal(body, &links)
	if err != nil {
		sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	err = storage.UpsertLinksWithGroups(ctx, GetLinkWithGroupGetSetters(links))
	switch {
	case err == nil:
		goto end
	case errors.Is(err, ErrFailedToUnmarshal):
		sendError(w, r, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrFailedUpsertLinks):
		// TODO: Break out errors into different status for different reasons
		fallthrough
	default:
		sendError(w, r, http.StatusInternalServerError, err.Error())
	}
end:
}

// GetLinkWithGroupGetSetters transforms a slice of LinkWithGroup into
// LinksWithGroupsGetSetter which is also a slice of LinkWithGroup.
func GetLinkWithGroupGetSetters(links []LinkWithGroup) storage.LinksWithGroupsGetSetter {
	ll := make([]storage.LinkWithGroupPropGetSetter, len(links))
	for i, link := range links {
		ll[i] = link
	}
	return linksWithGroups{links: ll}
}

var _ storage.LinksWithGroupsGetSetter = (*linksWithGroups)(nil)

type linksWithGroups struct {
	links []storage.LinkWithGroupPropGetSetter
}

func (ls linksWithGroups) GetLinkCount() int {
	return len(ls.links)
}

func (ls linksWithGroups) GetLinksWithGroups() []storage.LinkWithGroupPropGetSetter {
	return ls.links
}

func (ls linksWithGroups) SetLinksWithGroups(links []storage.LinkWithGroupPropGetSetter) {
	ls.links = links
}
