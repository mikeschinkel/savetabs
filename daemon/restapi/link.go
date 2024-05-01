package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"savetabs/sqlc"
	"savetabs/storage"
)

var _ storage.LinkGetSetter = (*link)(nil)

type link struct {
	Id          int64             `json:"id,omitempty"`
	OriginalUrl string            `json:"original_url,omitempty"`
	Title       string            `json:"title"`
	MetaMap     map[string]string `json:"meta,omitempty"`
	Scheme      string            `json:"scheme,omitempty"`
	Host        string            `json:"host,omitempty"`
	Subdomain   string            `json:"subdomain,omitempty"`
	Tld         string            `json:"tld,omitempty"`
	Sld         string            `json:"sld,omitempty"`
	Port        string            `json:"port,omitempty"`
	Path        string            `json:"path,omitempty"`
	Query       string            `json:"query,omitempty"`
	Fragment    string            `json:"fragment,omitempty"`
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

	if linkUrl != lnk.OriginalUrl {
		a.sendError(w, r,
			http.StatusBadRequest,
			fmt.Sprintf("Mismatch in URL in body vs. in URL path: bodyURL=%s, pathURL=%s",
				lnk.OriginalUrl,
				linkUrl,
			),
		)
		return
	}

	slog.Debug("PUT Link by URL:", "link", lnk)

	var linkId int64
	db = sqlc.GetNestedDBTX(sqlc.GetDatastore())
	err = db.Exec(func(dbtx sqlc.DBTX) error {
		linkId, err = storage.UpsertLink(ctx, dbtx.(*sqlc.NestedDBTX), &lnk)
		return err
	})
	switch {
	case err == nil:
		w.Header().Set("Location", LinkEndpoint(linkId))
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
	var link sqlc.LoadLinkRow
	var metas []sqlc.ListLinkMetaForLinkIdRow
	ctx := context.TODO()
	db := sqlc.GetNestedDBTX(sqlc.GetDatastore())
	err := db.Exec(func(dbtx sqlc.DBTX) (err error) {
		var q = a.DataStore.Queries(dbtx)

		link, err = q.LoadLink(ctx, linkId)
		if err != nil {
			goto end
		}
		metas, err = q.ListLinkMetaForLinkId(ctx, linkId)
		if err != nil {
			goto end
		}
	end:
		return err
	})
	if err != nil {
		a.sendError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, linkResponseFromLoadLinkRow(link, metas))
}

func linkResponseFromLoadLinkRow(row sqlc.LoadLinkRow, metas listLinkMetaForLinkIdRows) (lnk link) {
	return link{
		Id:          row.ID,
		OriginalUrl: row.OriginalUrl,
		Title:       row.Title,
		Scheme:      row.Scheme,
		Subdomain:   row.Subdomain,
		Tld:         row.Tld,
		Sld:         row.Sld,
		Port:        row.Port,
		Path:        row.Path,
		Query:       row.Query,
		Fragment:    row.Fragment,
		MetaMap:     metas.metaMap(),
	}
}

type listLinkMetaForLinkIdRows []sqlc.ListLinkMetaForLinkIdRow

func (l listLinkMetaForLinkIdRows) metaMap() map[string]string {
	mm := make(map[string]string, len(l))
	for _, row := range l {
		mm[row.Key] = row.Value
	}
	return mm
}
