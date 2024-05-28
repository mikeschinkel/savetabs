package ui

import (
	"embed"
	"fmt"

	"github.com/google/safehtml/template"
)

//go:embed html/*.template.html
var htmlFS embed.FS
var trustedHTMLFS = template.TrustedFSFromEmbed(htmlFS)

func GetTemplate(name string) *template.Template {
	name = fmt.Sprintf("%s.template.html", name)
	return template.Must(template.New(name).ParseFS(trustedHTMLFS, "html/"+name))
}

func init() {
	elements := template.MakeTrustedStringSlice("a", "li", "section", "img", "div", "expand-icon", "span", "form", "input")

	template.AddTrustedElementsAndAttributesForContext("url", elements,
		template.MakeTrustedStringSlice("hx-target", "hx-get", "hx-put", "hx-delete", "hx-post", "hx-patch"),
	)
	template.AddTrustedElementsAndAttributesForContext("identifier", elements,
		template.MakeTrustedStringSlice("hx-trigger", "hx-sync", "x-data", "x-init"),
	)
}
