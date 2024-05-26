package guard

import (
	"errors"
	"fmt"
	"net/http"

	"savetabs/model"
	"savetabs/shared"
	"savetabs/ui"
)

type LinksetToAdd struct {
	Action  string
	LinkIds []string
}

func AddLinksetIfNotExists(ctx Context, ls LinksetToAdd) (err error) {
	var linkIds []int64

	if !shared.ValidateAction(ls.Action) {
		err = errors.Join(ErrInvalidUpsertAction, fmt.Errorf("action=%s", ls.Action))
		goto end
	}
	linkIds, err = ParseLinkIds(ls.LinkIds)
	if err != nil {
		goto end
	}
	err = model.AddLinksetIfNotExist(ctx, model.LinksetToAdd{
		Action:  shared.NewAction(ls.Action),
		LinkIds: linkIds,
	})
	if err != nil {
		goto end
	}
end:
	return err
}

func GetLinksetSuccessAlertHTML(ctx Context, linkIds []string) (HTMLResponse, error) {
	hr := ui.NewHTMLResponse()
	ids, err := ParseLinkIds(linkIds)
	if err != nil {
		hr.SetCode(http.StatusInternalServerError)
		goto end
	}
	hr = ui.GetLinksetSuccessAlertHTML(ctx, ids)
end:
	return HTMLResponse{hr}, err
}

type LinksetParams struct {
	FilterQuery *shared.FilterQuery
	RequestURI  string
	Host        string
}

func GetLinksetHTML(ctx Context, p LinksetParams) (_ HTMLResponse, err error) {
	var fq shared.FilterQuery

	hr := ui.NewHTMLResponse()
	switch {
	case p.FilterQuery != nil:
		// This occurs when page is reloaded after a POST and the query was saved
		// serialized in a hidden field.
		fq = *p.FilterQuery
	default:
		fq, err = shared.ParseFilterQuery(p.RequestURI)
		if err != nil {
			goto end
		}
	}
	hr, err = ui.GetLinksetHTML(ctx, ui.LinksetArgs{
		FilterQuery: fq,
		RequestURI:  shared.MakeSafeURL(p.RequestURI),
		APIURL:      shared.MakeSafeAPIURL(p.Host),
	})
end:
	return HTMLResponse{hr}, err
}
