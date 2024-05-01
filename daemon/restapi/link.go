package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"savetabs/sqlc"
	"savetabs/storage"
)

var _ storage.LinkGetSetter = (*link)(nil)

type link struct {
	Id          int64             `json:"id"`
	OriginalUrl string            `json:"original_url"`
	Title       string            `json:"title"`
	MetaMap     map[string]string `json:"meta"`
}

func (l *link) GetId() int64 {
	return l.Id
}

func (l *link) SetId(id int64) {
	l.Id = id
}

func (l *link) GetTitle() string {
	return l.Title
}

func (l *link) SetTitle(t string) {
	l.Title = t
}

func (l *link) GetOriginalURL() string {
	return l.OriginalUrl
}

func (l *link) SetOriginalURL(u string) {
	l.OriginalUrl = u
}

func (l *link) GetMetaMap() map[string]string {
	return l.MetaMap
}

func (l *link) SetMetaMap(t map[string]string) {
	l.MetaMap = t
}

func (a *API) PutLinksByUrlLinkUrl(w http.ResponseWriter, r *http.Request, linkUrl LinkUrl) {
	var db *sqlc.NestedDBTX

	// TODO: Be sure to set location in HTTP header when link created
	ctx := context.TODO()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Find a better status result than "Bad Gateway"
		a.sendError(w, r, http.StatusBadGateway, err.Error())
		return
	}
	var lnk link
	err = json.Unmarshal(body, &lnk)
	if err != nil {
		a.sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	slog.Debug("PUT Link by URL:", "link", lnk)

	var linkId int64
	db = sqlc.GetNestedDBTX(sqlc.GetDatastore())
	err = db.Exec(func(sqlc.DBTX) error {
		linkId, err = storage.UpsertLink(ctx, db, &lnk)
		return err
	})
	switch {
	case err == nil:
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
	sendJSON(w, http.StatusOK, `{"id":0}`)
}
