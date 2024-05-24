package restapi

import (
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

func (f linkSetForm) FilterQuery() (fq shared.FilterQuery, err error) {
	print()
	//	var params GetHtmlLinksetParams
	//	var filterItems  []shared.FilterItem
	//
	//	err = json.Unmarshal([]byte(f.queryJSON), &params)
	//	if err != nil {
	//		goto end
	//	}
	//	filterItems = make([]shared.FilterItem,3)
	//	shared.ParseGroupTypeFilter()
	//	filterItems[0] = shared.GroupTypeFilter{GroupTypes: *params.Gt,}
	//
	//
	//}
	//	fq = shared.FilterQuery{FilterItems: []shared.FilterItem{
	//		*params.Gt,
	//
	//	}}
	//end:
	return fq, err
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
		fq, err := form.FilterQuery()
		if err != nil {
			a.sendError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		hr, err := guard.GetLinksetHTML(ctx, guard.LinksetParams{
			FilterQuery: &fq,
			RequestURI:  r.RequestURI,
			Host:        r.Host,
		})
		if err != nil {
			a.sendError(w, r, hr.Code(), err.Error())
			return
		}
		alert, err = guard.GetLinksetSuccessAlertHTML(ctx, form.LinkIds)
		if err != nil {
			//goland:noinspection GoDfaErrorMayBeNotNil
			a.sendError(w, r, alert.Code(), err.Error())
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
