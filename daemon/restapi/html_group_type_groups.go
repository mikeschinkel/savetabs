package restapi

import (
	"context"
	"net/http"

	"savetabs/ui"
)

func (a *API) GetHtmlGroupTypesGroupTypeNameGroups(w http.ResponseWriter, r *http.Request, groupTypeName GroupTypeName) {
	sendWith(context.Background(), w, r, func(ctx Context) ([]byte, error) {
		return ui.GroupTypeGroupsHTML(ctx, r.Host, groupTypeName)
	})
}
