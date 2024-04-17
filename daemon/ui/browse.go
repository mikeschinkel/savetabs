package ui

import (
	"bytes"
	"net/http"
)

var browseTemplate = GetTemplate("browse")

type browse struct {
	Host       string
	GroupTypes []groupType
}

func GetBrowseHTML(ctx Context, host string) (html []byte, status int, err error) {
	var out bytes.Buffer
	var gts []groupType

	gg, err := queries.ListGroupsType(ctx)
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
	return html, http.StatusInternalServerError, err
}
