package guard

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"savetabs/model"
	"savetabs/shared"
	"savetabs/ui"
)

type LinksetToAdd struct {
	Action    string
	LinkIds   []string
	QueryJSON string
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
	var hr ui.HTMLResponse
	ids, err := ParseLinkIds(linkIds)
	if err != nil {
		hr.HTTPStatus = http.StatusInternalServerError
		goto end
	}
	hr = ui.GetLinksetSuccessAlertHTML(ctx, ids)
end:
	return HTMLResponse(hr), err
}

type LinksetParams struct {
	// Gt Links for a Group Type
	GroupTypeFilter []string `form:"gt,omitempty" json:"gt,omitempty"`
	// Grp links by group
	GroupFilter []string `form:"grp,omitempty" json:"g,omitempty"`
	// M Key/Value meta filter for Links
	MetaFilter map[string]string `form:"m,omitempty" json:"m,omitempty"`
}

func GetLinksetHTML(ctx Context, host, requestURI string, p LinksetParams) (_ HTMLResponse, err error) {
	var hr ui.HTMLResponse

	rURI, err := url.Parse(requestURI)
	if err != nil {
		goto end
	}
	hr, err = ui.GetLinksetHTML(ctx, ui.LinksetParams{
		Host:       shared.NewHost(host),
		RequestURI: rURI,
		FilterQuery: shared.FilterQuery{
			FilterLabel: shared.NewLabel(p.getFilterLabel()),
			FilterTypes: shared.FilterTypes,
			Filters: shared.FilterMap{
				ui.GroupTypeFilter: shared.NewSliceFilter(ui.GroupTypeFilter, p.GroupTypeFilter),
				ui.GroupFilter:     shared.NewSliceFilter(ui.GroupFilter, p.GroupFilter),
				ui.MetaFilter:      shared.NewMapFilter(ui.MetaFilter, p.MetaFilter),
			},
		},
	})
end:
	return HTMLResponse(hr), err
}

func (lp LinksetParams) getFilterLabel() string {
	var name string
	sb := strings.Builder{}
	for _, ft := range shared.FilterTypes {
		switch ft {
		case GroupTypeFilter:
			name = "Group Type"
		case GroupFilter:
			name = "Group"
		case MetaFilter:
			name = "Meta"
		default:
			name = fmt.Sprintf("Unexpected[%s]", ft)
		}
		values := lp.getFilterValues(ft.String())
		if len(values) == 0 {
			continue
		}
		sb.WriteString(fmt.Sprintf("%s(s): %s && ", name, values))
	}
	if sb.Len() == 0 {
		return ""
	}
	labels := sb.String()
	// Strip off trailing ' && ' with -4
	return labels[:len(labels)-4]
}

func (lp LinksetParams) getFilterValues(ft string) (filters []string) {
	switch strings.ToUpper(ft) {
	case GroupTypeFilter.String():
		return shared.ToUpperSlice(lp.GroupTypeFilter)
	case GroupFilter.String():
		return lp.GroupFilter
	case MetaFilter.String():
		if lp.MetaFilter == nil {
			filters = []string{}
			goto end
		}
		filters = make([]string, len(lp.MetaFilter))
		i := 0
		for key, value := range lp.MetaFilter {
			filters[i] = fmt.Sprintf("%s=%s", key, value)
		}
	default:
		filters = []string{}
	}
end:
	return filters
}

//// TODO: Find where this is needed after refactor.
//func MetaFilterString(mf map[string]string) (s string) {
//	sb := strings.Builder{}
//	for k, v := range mf {
//		sb.WriteString(fmt.Sprintf("key[%s]=%s", k, v))
//		sb.WriteByte('&')
//	}
//	s = sb.String()
//	if len(s) == 0 {
//		return ""
//	}
//	return s[:len(s)-1]
//}
