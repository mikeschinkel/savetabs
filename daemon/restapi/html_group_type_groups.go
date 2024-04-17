package restapi

import (
	"context"
	"net/http"

	"savetabs/ui"
)

func (a *API) GetHtmlGroupTypesGroupTypeNameGroups(w http.ResponseWriter, r *http.Request, groupTypeName GroupTypeName) {
	sendWith(context.Background(), w, r, func(ctx Context) ([]byte, int, error) {
		return ui.GetGroupTypeGroupsHTML(ctx, r.Host, groupTypeName)
	})
}
