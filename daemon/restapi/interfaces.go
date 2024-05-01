package restapi

import (
	"net/http"

	"savetabs/ui"
)

type Viewer = ui.Viewer

type Persister interface {
	PostLinksWithGroups(http.ResponseWriter, *http.Request)
	PutLinksByUrlLinkUrl(http.ResponseWriter, *http.Request, LinkUrl)
	GetLinksLinkId(http.ResponseWriter, *http.Request, LinkId)
}
