package restapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"savetabs/guard"
	"savetabs/shared"
)

const (
	actionField    = "action"
	linkIdField    = "link_id"
	queryJSONField = "query_json"
)

type linkSetForm struct {
	Action    string
	LinkIds   []string
	queryJSON string
}

func (f linkSetForm) linksetParams() (lsp guard.LinksetParams, err error) {
	var params GetHtmlLinksetParams
	err = json.Unmarshal([]byte(f.queryJSON), &params)
	if err != nil {
		goto end
	}
	lsp = params.linksetParams()
end:
	return lsp, err
}

func (p GetHtmlLinksetParams) linksetParams() guard.LinksetParams {
	return guard.LinksetParams{
		GroupTypeFilter: shared.EnsureNotNil(p.Gt, []string{}),
		TabGroupFilter:  shared.EnsureNotNil(p.G, []string{}),
		CategoryFilter:  shared.EnsureNotNil(p.C, []string{}),
		TagFilter:       shared.EnsureNotNil(p.T, []string{}),
		KeywordFilter:   shared.EnsureNotNil(p.K, []string{}),
		BookmarkFilter:  shared.EnsureNotNil(p.B, []string{}),
		MetaFilter:      shared.EnsureNotNil(p.M, map[string]string{}),
	}
}

func parseLinksetForm(form url.Values) (lsf linkSetForm, err error) {
	var queryJSONs []string

	lsf.Action = form.Get(actionField)

	_, ok := form[linkIdField]
	if !ok {
		err = ErrNoLinkIdsSubmitted
		goto end
	}
	lsf.LinkIds = make([]string, len(form[linkIdField]))
	for i, id := range form[linkIdField] {
		lsf.LinkIds[i] = id
	}
	queryJSONs, ok = form[queryJSONField]
	if !ok {
		queryJSONs = []string{"{}"}
	}
	// TODO: Verify this is still valid given the Ardan Labs changes
	lsf.queryJSON = strings.Join(queryJSONs, "")

end:
	return lsf, err
}

func (a *API) PostHtmlLinkset(w http.ResponseWriter, r *http.Request) {
	var form linkSetForm

	err := r.ParseForm()
	if err != nil {
		a.sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	form, err = parseLinksetForm(r.Form)
	if err != nil {
		a.sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := newContext()
	err = guard.AddLinksetIfNotExists(ctx, guard.LinksetToAdd{
		Action:  form.Action,
		LinkIds: form.LinkIds,
	})

	switch {
	case err == nil:
		var alert guard.HTMLResponse
		lsp, err := form.linksetParams()
		if err != nil {
			a.sendError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		hr, err := guard.GetLinksetHTML(ctx, r.Host, r.RequestURI, lsp)
		if err != nil {
			a.sendError(w, r, hr.HTTPStatus, err.Error())
			return
		}
		alert, err = guard.GetLinksetSuccessAlertHTML(ctx, form.LinkIds)
		if err != nil {
			//goland:noinspection GoDfaErrorMayBeNotNil
			a.sendError(w, r, alert.HTTPStatus, err.Error())
			return
		}
		a.sendHTML(w, hr.HTML, alert.HTML)
		return
	case errors.Is(err, ErrFailedToUnmarshal):
		a.sendError(w, r, http.StatusBadRequest, err.Error())
		return
	case errors.Is(err, ErrFailedUpsertLinks):
		// TODO: Break out errors into different status for different reasons
		fallthrough
	default:
		a.sendError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}
