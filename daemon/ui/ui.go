package ui

import (
	"embed"
	"fmt"

	"github.com/google/safehtml/template"
	"savetabs/sqlc"
)

//go:embed html/*.template.html
var htmlFS embed.FS
var trustedHTMLFS = template.TrustedFSFromEmbed(htmlFS)

func GetTemplate(name string) *template.Template {
	name = fmt.Sprintf("%s.template.html", name)
	return template.Must(template.New(name).ParseFS(trustedHTMLFS, "html/"+name))
}

var dataStore sqlc.DataStore
var queries *sqlc.Queries

func Initialize(ctx Context, ds sqlc.DataStore) (err error) {
	dataStore = ds
	queries = ds.Queries()

	elements := template.MakeTrustedStringSlice("a", "li", "section", "img", "div")
	template.AddTrustedElementsAndAttributesForContext("url", elements,
		template.MakeTrustedStringSlice("hx-get"),
	)
	template.AddTrustedElementsAndAttributesForContext("identifier", elements,
		template.MakeTrustedStringSlice("hx-target", "hx-trigger", "hx-sync"),
	)

	return err
}
