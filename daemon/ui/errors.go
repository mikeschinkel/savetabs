package ui

import (
	"errors"
	"net/http"

	"savetabs/shared"
)

type HttpError struct {
	error
	StatusCode int
}

func NewHTTPError(code int, msg string) HttpError {
	return HttpError{
		error:      errors.New(msg),
		StatusCode: code,
	}
}

var errorTemplate = GetTemplate("error")

type ErrorParams struct {
	Err  error
	Host shared.Host
}

func GetErrorHTML(p ErrorParams) (hr HTMLResponse, err error) {
	var httpErr HttpError
	hr = NewHTMLResponse()
	hr.StatusCode = http.StatusInternalServerError

	if errors.As(err, &httpErr) {
		hr.StatusCode = httpErr.StatusCode
	}
	hr.HTML, err = errorTemplate.ExecuteToHTML(p.Err)
	if err != nil {
		goto end
	}
end:
	return hr, err
}
