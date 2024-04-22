package restapi

import (
	"net/http"

	"savetabs/ui"
)

type Viewer = ui.Viewer

type Persister interface {
	PostLinksWithGroups(http.ResponseWriter, *http.Request)
}
