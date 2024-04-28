package ui

import (
	"net/http"

	"github.com/google/safehtml"
)

type Message struct {
	Text  string
	Items []string
}

func (m *Message) HasItems() bool {
	return len(m.Items) > 0
}

type Alert struct {
	alertType  AlertType
	Message    Message
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
var alertOOBTemplate = GetTemplate("alert-oob")

func (*Views) GetAlertHTML(_ Context, typ AlertType, msg Message) (html safehtml.HTML, _ int, err error) {
	alert := &Alert{
		alertType:  typ,
		Message:    msg,
		HTTPStatus: http.StatusOK, // TODO: Might need more than this
	}
	html, err = alertTemplate.ExecuteToHTML(alert)
	if err != nil {
		goto end
	}
end:
	return html, alert.HTTPStatus, err
}

func (v *Views) GetOOBAlertHTML(ctx Context, typ AlertType, msg Message) (html safehtml.HTML, status int, err error) {
	html, status, err = v.GetAlertHTML(ctx, typ, msg)
	if err != nil {
		goto end
	}
	html, err = alertOOBTemplate.ExecuteToHTML(struct {
		AlertHTML safehtml.HTML
	}{
		AlertHTML: html,
	})
	if err != nil {
		goto end
	}
	if status == 0 {
		status = http.StatusOK
	}
end:
	return html, status, err
}
