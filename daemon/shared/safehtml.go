package shared

import (
	"fmt"
	"log/slog"

	"github.com/google/safehtml"
	"github.com/google/safehtml/template"
)

var ToSafeId = safehtml.IdentifierFromConstant
var ToPrefixedSafeId = safehtml.IdentifierFromConstantPrefix

var valueTemplate = template.Must(template.New("").Parse(`{{.Value}}`))

func MakeSafeHTML(h string) safehtml.HTML {
	sh, err := valueTemplate.ExecuteToHTML(struct{ Value string }{Value: h})
	if err != nil {
		panic(err)
	}
	return sh
}

func MakeSafeHTMLf(format string, args ...any) safehtml.HTML {
	return MakeSafeHTML(fmt.Sprintf(format, args...))
}

func MakeSafeURL(u string) safehtml.URL {
	return safehtml.URLSanitized(u)
}

func MakeSafeAPIURL(host string) safehtml.URL {
	return safehtml.URLSanitized(Host{host}.URL())
}

func MustMakeSafeJSON(j string) safehtml.JSON {
	jsn, err := safehtml.JSONFromValue(j)
	if err != nil {
		panic(err)
	}
	return jsn
}

func MakeSafeJSON(j string) (safehtml.JSON, error) {
	return safehtml.JSONFromValue(j)
}

func MakeEmptyObjectJSON() safehtml.JSON {
	return safehtml.EmptyObjectJSON()
}

func MakeSafeURLf(format string, args ...any) safehtml.URL {
	return MakeSafeURL(fmt.Sprintf(format, args...))
}

func MakeSafeId(id string) safehtml.Identifier {
	sid := safehtml.IdentifierSanitized(id)
	_id := sid.String()
	if _id[0] != 'i' {
		goto end
	}
	if _id[:8] == "invalid:" {
		slog.Error("HTML id found to be invalid", "html_id", id)
	}
end:
	return sid
}

func MakeSafeIdf(format string, args ...any) safehtml.Identifier {
	return MakeSafeId(fmt.Sprintf(format, args...))
}
