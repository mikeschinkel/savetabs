package restapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"savetabs/ui"
)

var _ ui.FilterGetter = (*GetHtmlLinksetParams)(nil)

func GetHtmlLinksetParamsFromJSON(j string) (ui.FilterGetter, error) {
	var p GetHtmlLinksetParams
	err := json.Unmarshal([]byte(j), &p)
	return p, err
}

func (p GetHtmlLinksetParams) GetFilterJSON() (string, error) {
	b, err := json.Marshal(p)
	return string(b), err
}

func (p GetHtmlLinksetParams) GetFilterLabels() string {
	var name string
	sb := strings.Builder{}
	for _, ft := range ui.FilterTypes {
		switch strings.ToUpper(ft) {
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
			name = fmt.Sprintf("Unexpected[%s]", ft)
		}
		values := p.GetFilterValues(ft)
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
			filters[i] = fmt.Sprintf("%s=%s", key, value)
		}
	default:
		filters = []string{}
	}
end:
	return filters
}

func (f MetadataFilter) String() (s string) {
	sb := strings.Builder{}
	for k, v := range f {
		sb.WriteString(fmt.Sprintf("key[%s]=%s", k, v))
		sb.WriteByte('&')
	}
	s = sb.String()
	if len(s) == 0 {
		return ""
	}
	return s[:len(s)-1]
}
