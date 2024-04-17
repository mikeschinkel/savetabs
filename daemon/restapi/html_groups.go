package restapi

import (
	"context"
	"net/http"

	"savetabs/ui"
)

func (a *API) GetHtmlGroupsGroupTypeGroupSlug(w http.ResponseWriter, r *http.Request, groupType GroupType, groupSlug GroupSlug) {
	sendWith(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return ui.GetGroupHTML(ctx, r.Host, groupType, groupSlug)
	})
}
