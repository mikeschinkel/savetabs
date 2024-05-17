package guard

import (
	"bytes"
	"log/slog"
	"strings"

	"golang.org/x/net/html"
	"savetabs/model"
)

func ParseContent(h string) (c model.Content, err error) {
	var doc *html.Node

	doc, err = html.Parse(strings.NewReader(h))
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
	c.Head, err = parseHTMLFragment("head", doc.FirstChild)
	if err != nil {
		goto end
	}
	c.Body, err = parseHTMLFragment("body", doc.LastChild)
	if err != nil {
		goto end
	}
end:
	if err != nil {
		slog.Warn("Unable to parse HTML",
			"html", h,
			"error", err,
		)
	}
	return c, err
}

func parseHTMLFragment(elem string, node *html.Node) (hf model.HTMLFragment, err error) {
	if node.Data == "head" {
		b := bytes.Buffer{}
		err = html.Render(&b, node)
		if err != nil {
			goto end
		}
		hf, err = model.ParseHTMLFragment(b.String())
		if err != nil {
			goto end
		}
	}
end:
	return hf, err
}
