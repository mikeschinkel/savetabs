package guard

import (
	"errors"

	"github.com/mikeschinkel/savetabs/daemon/shared"
	"github.com/mikeschinkel/savetabs/daemon/ui"
)

type ErrorParams struct {
	Host string
	Msg  string
}

func GetErrorHTML(params ErrorParams) (_ HTMLResponse, err error) {
	var hr ui.HTMLResponse
	hr, err = ui.GetErrorHTML(ui.ErrorParams{
		Err:  errors.New(params.Msg),
		Host: shared.NewHost(params.Host),
	})
	return HTMLResponse{hr}, err
}
