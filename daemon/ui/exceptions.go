package ui

import (
	"net/http"
)

type ExceptionsParams struct {
	Exceptions []string
}

func GetExceptionsHTML(ctx Context, p ExceptionsParams) (hr HTMLResponse, err error) {
	hr, err = GetAlertHTML(ctx, AlertParams{
		//Host: shared.Host{},
		OOB:  false,
		Type: AlertAlert,
		Message: Message{
			Text:  "Notice:",
			Items: p.Exceptions,
		},
	})
	if err != nil {
		hr.StatusCode = http.StatusInternalServerError
		goto end
	}
	hr.StatusCode = http.StatusAccepted
end:
	return hr, err
}
