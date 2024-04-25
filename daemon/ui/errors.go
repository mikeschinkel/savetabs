package ui

import (
	"bytes"
	"errors"
	"net/http"
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

func (*Views) GetErrorHTML(_ Context, err error) (html []byte, _ int, _ error) {
	var out bytes.Buffer
	var httpErr HttpError
	statusCode := http.StatusInternalServerError

	if errors.As(err, &httpErr) {
		statusCode = httpErr.StatusCode
	}
	err = errorTemplate.Execute(&out, err)
	if err != nil {
		goto end
	}
	html = out.Bytes()
end:
	return html, statusCode, err
}
