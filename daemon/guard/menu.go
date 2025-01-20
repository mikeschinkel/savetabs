package guard

import (
	"github.com/mikeschinkel/savetabs/daemon/shared"
	"github.com/mikeschinkel/savetabs/daemon/ui"
)

type MenuParams struct {
	Host string
}

func GetMenuHTML(ctx Context, params MenuParams) (HTMLResponse, error) {
	hr, err := ui.GetMenuHTML(ctx, ui.HTMLMenuParams{
		Host: shared.NewHost(params.Host),
	})
	return HTMLResponse{hr}, err
}
