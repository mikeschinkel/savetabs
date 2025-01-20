package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mikeschinkel/savetabs/daemon/guard"
	"github.com/mikeschinkel/savetabs/daemon/shared"
)

type link struct {
	TabId int64  `json:"tab_id"`
	URL   string `json:"url"`
	Title string `json:"title"`
	HTML  string `json:"html"`
}

func (a *API) PutLinksByUrlLinkUrl(w http.ResponseWriter, r *http.Request, linkUrl LinkUrl) {
	ctx := context.TODO()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Find a better status result than "Bad Gateway"
		a.sendHTMLError(w, r, http.StatusBadGateway, err.Error())
		return
	}
	var link link
	err = json.Unmarshal(body, &link)
	if err != nil {
		a.sendHTMLError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if linkUrl != link.URL {
		a.sendHTMLError(w, r,
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
	linkId, err = guard.UpsertLink(ctx, guard.UpsertLinkArgs{
		URL:   link.URL,
		Title: link.Title,
		HTML:  link.HTML,
	})
	switch {
	case err == nil:
		w.Header().Set("Location", shared.LinkEndpointURL(linkId))
		a.sendHTMLStatus(w, http.StatusNoContent)
	default:
		a.sendHTMLError(w, r, http.StatusInternalServerError, err.Error())
	}
}

func (a *API) GetLinksLinkId(w http.ResponseWriter, r *http.Request, linkId LinkId) {
	link, err := guard.LoadLink(context.TODO(), linkId)
	if err != nil {
		a.sendHTMLError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, link)
}
