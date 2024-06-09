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
		// This occurs when linkset is reloaded after a POST and the query was saved
		// serialized in a hidden field.
		fq = *p.FilterQuery
	default:
		// This occurs when linkset is loaded directly from a URL.
		fq, err = shared.ParseFilterQuery(p.RequestURI)
		if err != nil {
			goto end
		}
	}
	err = updateFilterQueryParentDBId(ctx, &fq)
	if err != nil {
		goto end
	}
	hr, err = ui.GetLinksetHTML(ctx, ui.LinksetArgs{
		FilterQuery: fq,
		RequestURI:  shared.MakeSafeURL(p.RequestURI),
		APIURL:      shared.MakeSafeAPIURL(p.Host),
	})
end:
	return HTMLResponse{hr}, err
}

// filterQueryParentDBId gets the Drag parent DBId.
// NOTE: This is hardcoded to support only group.
//
//	It will need to be rewritten to be generic
//	when and if we have more complex queries.
func updateFilterQueryParentDBId(ctx Context, fq *shared.FilterQuery) (err error) {
	var id int64
	var g shared.FilterGroup
	if fq.ParentDBId != 0 {
		id = fq.ParentDBId
	}
	fi, index := fq.FilterItemByType(shared.GroupFilterType)
	gf, ok := fi.(shared.GroupFilter)
	if !ok {
		// Some filter queries may not have a group filter
		goto end
	}
	// Try to get the group id from the filter group
	id = gf.DBId()
	if id != 0 {
		goto end
	}
	if len(gf.FilterGroups) == 0 {
		shared.Panicf("FilterQuery has a FilterItem of GroupFilterType but that item has not FilterGroups: %v", fq)
	}
	// If filter group did not have the id then let's try to load if from the first
	// in the list of FilterGroups given its slug.
	g = gf.FilterGroups[0]
	id, err = model.LoadGroupIdBySlug(ctx, g.Slug())
	fq.FilterItems[index] = gf.WithDBId(id)
	fq.ParentDBId = id
end:
	return err
}
