package ui

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/google/safehtml"
)

type Context = context.Context

type HTMLResponse struct {
	HTML       safehtml.HTML
	code       int
	HTTPStatus int
}

func NewHTMLResponse() HTMLResponse {
	return HTMLResponse{
		code: http.StatusOK,
	}
}

func (hr HTMLResponse) Code() int {
	if hr.code == 0 {
		slog.Warn("HTTPStatus request without code being set",
			"callstack", string(debug.Stack()),
		)
		hr.code = http.StatusOK
	}
	return hr.code
}
func (hr HTMLResponse) SetCode(code int) {
	hr.code = code
}
