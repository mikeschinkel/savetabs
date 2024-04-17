package ui

import (
	"embed"
	"fmt"

	"github.com/google/safehtml/template"
	"savetabs/sqlc"
)

var _ Viewer = (*Views)(nil)

type Views struct {
	DataStore sqlc.DataStore
	Queries   *sqlc.Queries
}

func NewViews(ds sqlc.DataStore) *Views {
	v := &Views{
		DataStore: ds,
	}
	v.Queries = ds.Queries()
	return v
}

//go:embed html/*.template.html
var htmlFS embed.FS
var trustedHTMLFS = template.TrustedFSFromEmbed(htmlFS)

func GetTemplate(name string) *template.Template {
	name = fmt.Sprintf("%s.template.html", name)
	return template.Must(template.New(name).ParseFS(trustedHTMLFS, "html/"+name))
}

func init() {
	elements := template.MakeTrustedStringSlice("a", "li", "section", "img", "div", "expand-icon", "span")

	template.AddTrustedElementsAndAttributesForContext("url", elements,
		template.MakeTrustedStringSlice("hx-get"),
	)
	template.AddTrustedElementsAndAttributesForContext("identifier", elements,
		template.MakeTrustedStringSlice("hx-target", "hx-trigger", "hx-sync", "x-data", "x-init"),
	)
}
