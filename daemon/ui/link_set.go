package ui

import (
	"fmt"
	"strconv"

	"github.com/google/safehtml"
)

type linkSet struct {
	apiURL string
	Links  []link
}

func (ls linkSet) HeaderHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`links-row-0`)
}
func (ls linkSet) FooterHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstantPrefix(`links-row`,
		strconv.Itoa(ls.NumLinks()),
	)
}
func (ls linkSet) NumLinks() int {
	return len(ls.Links)
}
func (ls linkSet) HTMLLinksURL() string {
	return fmt.Sprintf("%s/html/linkset", ls.apiURL)
}
