package guard

import (
	"errors"

	"savetabs/shared"
	"savetabs/ui"
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
	return HTMLResponse(hr), err
}
