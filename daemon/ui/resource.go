package ui

import (
	"html"
	"strconv"

	"github.com/google/safehtml"
	"savetabs/sqlc"
)

var resourcesTemplate = GetTemplate("resources")

type resource struct {
	Id        int64
	URL       string
	Domain    string
	Host      string
	GroupName string
	GroupType string
}

func (r resource) Title() string {
	return html.EscapeString(r.URL)
}
func (r resource) ARIALabel() string {
	return "External Link: " + html.EscapeString(r.URL)
}

func (r resource) Identifier() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`resource`,
		strconv.FormatInt(r.Id, 10),
	)
}

func constructResources(rfgs []sqlc.ListResourcesForGroupRow) []resource {
	rr := make([]resource, len(rfgs))
	for i, rfg := range rfgs {
		r := &resource{
			Id:        rfg.ID.Int64,
			URL:       rfg.Url.String,
			GroupName: rfg.GroupName,
			GroupType: rfg.GroupType,
			Domain:    rfg.Domain.String,
		}
		rr[i] = *r
	}
	return rr
}
