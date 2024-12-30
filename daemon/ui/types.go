package ui

import (
	"context"
	"net/http"

	"github.com/google/safehtml"
)

type Context = context.Context

type HTMLResponse struct {
	HTML       safehtml.HTML
	StatusCode int
}

func NewHTMLResponse() HTMLResponse {
	return HTMLResponse{
		StatusCode: http.StatusOK,
	}
}
