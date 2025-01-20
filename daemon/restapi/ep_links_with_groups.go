package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/mikeschinkel/savetabs/daemon/guard"
	"github.com/mikeschinkel/savetabs/daemon/shared"
)

type linkWithGroupForJSON struct {
	URL       string `json:"url"`
	Title     string `json:"title"`
	GroupId   int64  `json:"groupId"`
	GroupType string `json:"groupType"`
	Group     string `json:"group"`
}

func (a *API) PostLinksWithGroups(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Find a better status result than "Bad Gateway"
		a.sendHTMLError(w, r, http.StatusBadGateway, err.Error())
		return
	}
	var jsonLinks []linkWithGroupForJSON
	err = json.Unmarshal(body, &jsonLinks)
	if err != nil {
		a.sendHTMLError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = guard.AddLinksWithGroupsIfNotExists(ctx, guard.AddLinksWithGroupsParams{
		Links: shared.ConvertSlice(jsonLinks, func(link linkWithGroupForJSON) guard.LinkWithGroup {
			return guard.LinkWithGroup{
				URL:       link.URL,
				Title:     link.Title,
				GroupId:   link.GroupId,
				GroupType: link.GroupType,
				Group:     link.Group,
			}
		}),
	})

	switch {
	case err == nil:
		sendJSON(w, http.StatusOK, newJSONResponse(true))
	case errors.Is(err, ErrUnmarshallingJSON):
		sendJSON(w, http.StatusBadRequest, newJSONResponse(false))
	case errors.Is(err, ErrFailedUpsertLinks):
		// TODO: Break out errors for all different error types
		fallthrough
	default:
		sendJSON(w, http.StatusInternalServerError, newJSONResponse(false))
	}
}
