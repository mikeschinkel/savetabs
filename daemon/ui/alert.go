package ui

import (
	"net/http"

	"github.com/google/safehtml"
	"savetabs/shared"
)

type Message struct {
	Text  string
	Items []string
}

func NewMessage(text string, items []string) Message {
	return Message{Text: text, Items: items}
}

func (m *Message) HasItems() bool {
	return len(m.Items) > 0
}

type Alert struct {
	alertType AlertType
	Message   Message
}

type AlertType struct {
	value string
}

func NewAlertType(value string) AlertType {
	return AlertType{value: value}
}

var (
	UnspecifiedAlert = NewAlertType("")
	AlertAlert       = NewAlertType("alert")
	InfoAlert        = NewAlertType("info")
	SuccessAlert     = NewAlertType("success")
	WarningAlert     = NewAlertType("warning")
	ErrorAlert       = NewAlertType("error")
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

type AlertParams struct {
	Host    shared.Host
	OOB     bool
	Type    AlertType
	Message Message
}

type alertOOB struct {
	AlertHTML safehtml.HTML
}

func GetAlertHTML(_ Context, p AlertParams) (hr HTMLResponse, err error) {
	var html safehtml.HTML
	hr = NewHTMLResponse()
	alert := &Alert{
		alertType: p.Type,
		Message:   p.Message,
	}
	html, err = alertTemplate.ExecuteToHTML(alert)
	if err != nil {
		hr.StatusCode = http.StatusInternalServerError
		goto end
	}
	if !p.OOB {
		hr.HTML = html
		goto end
	}
	hr.HTML, err = alertOOBTemplate.ExecuteToHTML(alertOOB{
		AlertHTML: html,
	})
	if err != nil {
		hr.StatusCode = http.StatusInternalServerError
		goto end
	}
end:
	return hr, err
}
