package guard

import (
	"context"

	"github.com/mikeschinkel/savetabs/daemon/ui"
)

type Context = context.Context

type HTMLResponse struct {
	ui.HTMLResponse
}

func NewHTMLResponse() HTMLResponse {
	return HTMLResponse{ui.NewHTMLResponse()}
}
