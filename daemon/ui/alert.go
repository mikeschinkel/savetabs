package ui

import (
	"bytes"
	"net/http"

	"github.com/google/safehtml"
)

type Alert struct {
	alertType  AlertType
	Message    string
	HTTPStatus int
}

type AlertType string

const (
	UnspecifiedAlert AlertType = ""
	AlertAlert       AlertType = "alert"
	InfoAlert        AlertType = "info"
	SuccessAlert     AlertType = "success"
	WarningAlert     AlertType = "warning"
	ErrorAlert       AlertType = "error"
)

func NewAlert(typ AlertType, msg string, code int) *Alert {
	return &Alert{
		alertType:  typ,
		Message:    msg,
		HTTPStatus: code,
	}
}

func (a *Alert) IconHTML() safehtml.HTML {
	var name safehtml.HTML
	begin := safehtml.HTMLFromConstant("<")
	mid := safehtml.HTMLFromConstant(`-icon class="mx-1.5 mr-3"></`)
	end := safehtml.HTMLFromConstant(`-icon>`)
	switch a.alertType {
	case AlertAlert:
		name = safehtml.HTMLFromConstant(`alert`)
	case InfoAlert:
		name = safehtml.HTMLFromConstant(`info`)
	case SuccessAlert:
		name = safehtml.HTMLFromConstant(`success`)
	case WarningAlert:
		name = safehtml.HTMLFromConstant(`warning`)
	case ErrorAlert:
		name = safehtml.HTMLFromConstant(`error`)
	case UnspecifiedAlert:
		panic("AlertType is unspecified")
	default:
		panicf("AlertType is invalid: '%s'", a.alertType)
	}
	return safehtml.HTMLConcat(begin, name, mid, name, end)
}

func (a *Alert) AlertType() (id safehtml.Identifier) {
	switch a.alertType {
	case AlertAlert:
		id = safehtml.IdentifierFromConstant(`alert`)
	case InfoAlert:
		id = safehtml.IdentifierFromConstant(`info`)
	case SuccessAlert:
		id = safehtml.IdentifierFromConstant(`success`)
	case WarningAlert:
		id = safehtml.IdentifierFromConstant(`warning`)
	case ErrorAlert:
		id = safehtml.IdentifierFromConstant(`error`)
	case UnspecifiedAlert:
		panic("AlertType is unspecified")
	default:
		panicf("AlertType is invalid: '%s'", a.alertType)
	}
	return id
}

var alertTemplate = GetTemplate("alert")

func (*Views) GetAlertHTML(_ Context, typ AlertType, msg string) (html []byte, _ int, err error) {
	var out bytes.Buffer
	alert := &Alert{
		alertType:  typ,
		Message:    msg,
		HTTPStatus: http.StatusOK, // TODO: Might need more than this
	}
	err = alertTemplate.Execute(&out, alert)
	if err != nil {
		goto end
	}
	html = out.Bytes()
end:
	return html, alert.HTTPStatus, err
}
