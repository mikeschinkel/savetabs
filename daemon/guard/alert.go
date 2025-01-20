package guard

import (
	"github.com/mikeschinkel/savetabs/daemon/shared"
	"github.com/mikeschinkel/savetabs/daemon/ui"
)

type AlertParams struct {
	OOB   bool
	Host  string
	Type  string
	Msg   string
	Items []string
}

func GetAlertHTML(ctx Context, params AlertParams) (_ HTMLResponse, err error) {
	var hr ui.HTMLResponse
	hr, err = ui.GetAlertHTML(ctx, ui.AlertParams{
		OOB:     params.OOB,
		Host:    shared.NewHost(params.Host),
		Type:    ui.NewAlertType(params.Type),
		Message: ui.NewMessage(params.Msg, params.Items),
	})
	return HTMLResponse{hr}, err
}
