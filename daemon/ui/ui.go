package ui

import (
	"embed"
	"fmt"
	"sync"

	"github.com/google/safehtml/template"
	"savetabs/sqlc"
)

//go:embed html/*.template.html
var htmlFS embed.FS
var trustedHTMLFS = template.TrustedFSFromEmbed(htmlFS)

func getTemplate(name string) *template.Template {
	name = fmt.Sprintf("%s.template.html", name)
	return template.Must(template.New(name).ParseFS(trustedHTMLFS, "html/"+name))
}

var mutex sync.Mutex

var dataStore sqlc.DataStore
var queries *sqlc.Queries

func Initialize(ctx Context, ds sqlc.DataStore) (err error) {
	dataStore = ds
	queries = ds.Queries()

	elements := template.MakeTrustedStringSlice("a", "li", "section")
	template.AddTrustedElementsAndAttributesForContext("url", elements,
		template.MakeTrustedStringSlice("hx-get"),
	)
	template.AddTrustedElementsAndAttributesForContext("identifier", elements,
		template.MakeTrustedStringSlice("hx-target", "hx-trigger"),
	)

	return err
}
