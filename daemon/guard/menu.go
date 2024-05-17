package guard

import (
	"savetabs/shared"
	"savetabs/ui"
)

type MenuParams struct {
	Host string
}

func GetMenuHTML(ctx Context, params MenuParams) (HTMLResponse, error) {
	hr, err := ui.GetMenuHTML(ctx, ui.HTMLMenuParams{
		Host: shared.NewHost(params.Host),
	})
	return HTMLResponse(hr), err
}
