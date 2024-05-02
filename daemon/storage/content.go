package storage

import (
	"bytes"
	"log/slog"
	"strings"

	"golang.org/x/net/html"
)

type Content struct {
	LinkId int64
	Head   string
	Body   string
}

func (c *Content) setRawContent(h string) {
	doc, err := html.Parse(strings.NewReader(h))
	if err != nil {
		goto end
	}
	if doc == nil {
		err = ErrHTMLNotParsed
		goto end
	}
	if doc.Data == "" {
		doc = doc.FirstChild
	}
	if doc.FirstChild.Data == "head" {
		b := bytes.Buffer{}
		err = html.Render(&b, doc.FirstChild)
		if err != nil {
			goto end
		}
		c.Head = b.String()
	}
	if doc.LastChild.Data == "body" {
		b := bytes.Buffer{}
		err = html.Render(&b, doc.LastChild)
		if err != nil {
			goto end
		}
		c.Body = b.String()
	}
end:
	if err == nil {
		return
	}
	slog.Warn("Unable to parse HTML for link",
		"link_id", c.LinkId,
		"html", h,
		"error", err,
	)
}
