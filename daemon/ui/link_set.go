package ui

import (
	"fmt"

	"github.com/google/safehtml"
)

type linkSet struct {
	apiURL   string
	rawQuery string
	Label    string
	Links    []link
}

func (ls linkSet) HeaderHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`links-row-head`)
}
func (ls linkSet) FooterHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`links-row-foot`)
}
func (ls linkSet) URLQuery() safehtml.URL {
	return safehtml.URLSanitized("?" + ls.rawQuery)
}
func (ls linkSet) NumLinks() int {
	return len(ls.Links)
}
func (ls linkSet) HTMLLinksURL() string {
	return fmt.Sprintf("%s/html/linkset", ls.apiURL)
}
func (ls linkSet) TableHeaderFooterHTML() safehtml.HTML {
	return safehtml.HTMLFromConstant(`
<th class="p-0.5">#</th>
<th class="p-0.5">
	<label>
		<input type="checkbox" @change="maybeConfirmCheckAll" class="check-all"> 
	</label>
</th>
<th class="p-0.5 text-center">Link</th>
<th class="p-0.5 text-right">Sub</th>
<th class="p-0.5">Domain</th>
<th class="p-0.5">Title</th>
`)
}
