package restapi

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"savetabs/storage"
)

func (a *API) PostHtmlLinkset(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	err := r.ParseForm()
	if err != nil {
		sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	linkIds, ok := r.Form["link_id"]
	if !ok {
		sendError(w, r, http.StatusBadRequest, ErrNoLinkIdsSubmitted.Error())
		return
	}
	err = storage.UpsertLinkSet(ctx, linkSetAction{
		Action:  r.Form.Get("action"),
		LinkIds: linkIds,
	})
	switch {
	case err == nil:
		goto end
	case errors.Is(err, ErrFailedToUnmarshal):
		sendError(w, r, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrFailedUpsertLinks):
		// TODO: Break out errors into different status for different reasons
		fallthrough
	default:
		sendError(w, r, http.StatusInternalServerError, err.Error())
	}
end:
}

var _ storage.LinkSetActionGetter = (*linkSetAction)(nil)

type linkSetAction struct {
	Action  string
	LinkIds []string
}

func (l linkSetAction) GetAction() string {
	return l.Action
}

func (a linkSetAction) GetLinkIds() (ids []int64, err error) {
	var linkId int64

	ids = make([]int64, len(a.LinkIds))
	for i, id := range a.LinkIds {
		linkId, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			goto end
		}
		ids[i] = linkId
	}
end:
	return ids, err
}
