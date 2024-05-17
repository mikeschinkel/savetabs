package ui

import (
	"context"

	"github.com/google/safehtml"
)

type Context = context.Context

type HTMLResponse struct {
	HTML       safehtml.HTML
	HTTPStatus int
}
