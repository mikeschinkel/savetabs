package ui

import (
	"errors"
	"net/http"

	"github.com/google/safehtml"
)

var (
	ErrInvalidKeyFormat = errors.New("invalid key format (expected '<type>-<key>')")
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

func (*Views) GetErrorHTML(_ Context, err error) (html safehtml.HTML, _ int, _ error) {
	var httpErr HttpError
	statusCode := http.StatusInternalServerError

	if errors.As(err, &httpErr) {
		statusCode = httpErr.StatusCode
	}
	html, err = errorTemplate.ExecuteToHTML(err)
	if err != nil {
		goto end
	}
end:
	return html, statusCode, err
}
