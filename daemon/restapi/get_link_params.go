package restapi

import (
	"fmt"
	"strings"

	"savetabs/ui"
)

var _ ui.FilterValueGetter = (*GetHtmlLinksetParams)(nil)

func (p GetHtmlLinksetParams) RawQuery() string {
	//TODO implement me
	panic("implement me")
}

func (p GetHtmlLinksetParams) GetFilterLabel(typ, value string) string {
	var name string
	switch strings.ToUpper(typ) {
	case "GT":
		name = "Group Type"
	case "B":
		name = "Bookmark"
	case "C":
		name = "Categories"
	case "G":
		name = "Tab Group"
	case "K":
		name = "Keyword"
	case "T":
		name = "Tag"
	case "M":
		name = "Meta"
	default:
		name = fmt.Sprintf("Unexpected[%s]", typ)
	}
	// TODO: Remove (s) when only one value.
	//       That means `value string` needs to be `values []string`
	return fmt.Sprintf("%s(s): %s", name, value)
}

func (p GetHtmlLinksetParams) GetFilterValues(typ string) (filters []string) {
	ensureNotNil := func(ss *[]string) []string {
		if ss == nil {
			return []string{}
		}
		return *ss
	}
	switch strings.ToUpper(typ) {
	case "GT":
		return toUpperSlice(ensureNotNil(p.Gt))
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
			filters[i] = fmt.Sprintf("%typ=%typ", key, value)
		}
	default:
		filters = []string{}
	}
end:
	return filters
}
