package ui

import (
	"bytes"
	"context"
)

var browseTemplate = GetTemplate("browse")

type browse struct {
	Host       string
	GroupTypes []groupType
}

func BrowseHTML(host string) (html []byte, err error) {
	var out bytes.Buffer
	var gts []groupType

	gg, err := queries.ListGroupsType(context.Background())
	if err != nil {
		goto end
	}
	gts = newGroupTypeMap(gg).AsSortedSlice()
	err = browseTemplate.Execute(&out, browse{
		Host:       makeURL(host),
		GroupTypes: gts,
	})
	if err != nil {
		goto end
	}
	html = out.Bytes()
end:
	return html, err
}
