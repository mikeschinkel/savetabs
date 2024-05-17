package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"savetabs/guard"
	"savetabs/shared"
)

func (a *API) PutLinksByUrlLinkUrl(w http.ResponseWriter, r *http.Request, linkUrl LinkUrl) {
	ctx := context.TODO()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Find a better status result than "Bad Gateway"
		a.sendError(w, r, http.StatusBadGateway, err.Error())
		return
	}
	var link struct {
		TabId int64  `json:"tab_id"`
		URL   string `json:"url"`
		Title string `json:"title"`
		HTML  string `json:"html"`
	}
	err = json.Unmarshal(body, &link)
	if err != nil {
		a.sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if linkUrl != link.URL {
		a.sendError(w, r,
			http.StatusBadRequest,
			fmt.Sprintf("Mismatch in URL in body vs. in URL path: bodyURL=%s, pathURL=%s",
				link.URL,
				linkUrl,
			),
		)
		return
	}

	slog.Debug("PUT Link by URL:", "link", link)

	var linkId int64
	linkId, err = guard.LinkUpsert(ctx, guard.UpsertLink{
		URL:   link.URL,
		Title: link.Title,
		HTML:  link.HTML,
	})
	switch {
	case err == nil:
		w.Header().Set("Location", shared.LinkEndpointURL(linkId))
		sendJSON(w, http.StatusOK, struct {
			Id int64 `json:"id"`
		}{Id: linkId})
	case errors.Is(err, ErrFailedToUnmarshal):
		a.sendError(w, r, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrFailedUpsertLinks):
		// TODO: Break out errors into different status for different reasons
		fallthrough
	default:
		a.sendError(w, r, http.StatusInternalServerError, err.Error())
	}
}

func (a *API) GetLinksLinkId(w http.ResponseWriter, r *http.Request, linkId LinkId) {
	link, err := guard.LoadLink(context.TODO(), linkId)
	if err != nil {
		a.sendError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, link)
}
