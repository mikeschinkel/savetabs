package restapi

import (
	"context"
	"fmt"
	"net/http"

	"savetabs/ui"
)

func (a *API) GetLinks(w http.ResponseWriter, r *http.Request, params GetLinksParams) {
	sendWith(context.Background(), w, r, func(ctx Context) ([]byte, error) {
		return ui.GetLinksHTML(ctx, r.Host, params)
	})
}

var _ ui.SlugsForGetter = (*GetLinksParams)(nil)

func (p GetLinksParams) GetSlugsFor(gt GroupType) (filters []string) {
	ensureNotNil := func(ss *[]string) []string {
		if ss == nil {
			return []string{}
		}
		return *ss
	}
	switch gt {
	case "B":
		return ensureNotNil(p.B)
	case "C":
		return ensureNotNil(p.C)
	case "G":
		return ensureNotNil(p.G)
	case "K":
		return ensureNotNil(p.K)
	case "T":
		return ensureNotNil(p.T)
	case "M":
		if *p.M == nil {
			filters = []string{}
			goto end
		}
		filters = make([]string, len(*p.M))
		i := 0
		for key, value := range *p.M {
			filters[i] = fmt.Sprintf("%s=%s", key, value)
		}
	default:
		filters = []string{}
	}
end:
	return filters
}
