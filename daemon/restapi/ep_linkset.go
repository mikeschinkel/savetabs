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
	actionField      = "action"
	linkIdField      = "link_id"
	filterQueryField = "filter_query"
)

type linkSetForm struct {
	Action      string
	LinkIds     []string
	filterQuery string
}

func (f linkSetForm) FilterQuery() (fq shared.FilterQuery, err error) {
	return shared.ParseFilterQuery(f.filterQuery)
}

func parseLinksetForm(form url.Values) (lsf linkSetForm, err error) {
	var filterQueries []string

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
	filterQueries, ok = form[filterQueryField]
	if !ok {
		goto end
	}
	lsf.filterQuery = strings.Join(filterQueries, "&")

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
