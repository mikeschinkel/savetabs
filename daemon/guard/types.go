package guard

import (
	"context"

	"savetabs/ui"
)

type Context = context.Context

type HTMLResponse struct {
	ui.HTMLResponse
}

func NewHTMLResponse() HTMLResponse {
	return HTMLResponse{ui.NewHTMLResponse()}
}
