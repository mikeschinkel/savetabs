package ui

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/safehtml"
	"savetabs/model"
)

type htmlLinkset struct {
	apiURL     string
	Links      []model.Link
	Label      string
	requestURI string
	queryJSON  string
}

func (ls htmlLinkset) HTMLLinks() (links []htmlLink) {
	links = make([]htmlLink, len(ls.Links))
	for i, link := range ls.Links {
		links[i] = newHTMLLink(link, i+1)
	}
	return links
}
func (ls htmlLinkset) HeaderHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`links-row-head`)
}
func (ls htmlLinkset) FooterHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`links-row-foot`)
}
func (ls htmlLinkset) URLQuery() safehtml.URL {
	parts := strings.Split(ls.requestURI+"?", "?")
	return safehtml.URLSanitized("?" + parts[1])
}
func (ls htmlLinkset) QueryJSON() safehtml.JSON {
	j, err := safehtml.JSONFromValue(ls.queryJSON)
	if err != nil {
		slog.Error("Unable to create safe JSON", "json", ls.queryJSON)
		j = safehtml.JSONFromConstant("{}")
	}
	return j
}
func (ls htmlLinkset) NumLinks() int {
	return len(ls.Links)
}
func (ls htmlLinkset) HTMLLinksURL() string {
	return fmt.Sprintf("%s/html/linkset", ls.apiURL)
}
func (ls htmlLinkset) TableHeaderFooterHTML() safehtml.HTML {
	return safehtml.HTMLFromConstant(`
<th class="p-0.5">#</th>
<th class="p-0.5">
	<label>
		<input type="checkbox" @change="maybeConfirmCheckAll" class="check-all"> 
	</label>
</th>
<th class="p-0.5 text-center">Link</th>
<th class="p-0.5 max-w-[10vw]">Domain</th>
<th class="p-0.5 max-w-[15vw]">/Path</th>
<th class="p-0.5 max-w-[20vw]">?Query</th>
<th class="p-0.5 max-w-[20vw]">#Fragment</th>
<th class="p-0.5 max-w-[20vw]">Title</th>
<th class="p-0.5">Scheme</th>
<th class="p-0.5 text-right">Sub</th>
<th class="p-0.5">SLD</th>
<th class="p-0.5">TLD</th>
<th class="p-0.5">:Port</th>
`)
}

var linksetTemplate = GetTemplate("link-set")

type LinksetParams model.LinksetToLoadParams

func (lp LinksetParams) FilterJSON() (j string) {
	b, err := json.Marshal(lp.FilterQuery)
	if err != nil {
		j = "{}"
		slog.Error("Failed to marshal FilterQuery",
			"filter_query", lp.FilterQuery,
			"err", err.Error(),
		)
		goto end
	}
	j = string(b)
end:
	return j
}

func GetLinksetHTML(ctx Context, params LinksetParams) (hr HTMLResponse, err error) {
	var ls model.LinksetToLoad

	hr.HTTPStatus = http.StatusOK

	ls, err = model.LoadLinkset(ctx, model.LinksetToLoadParams(params))
	if err != nil {
		hr.HTTPStatus = http.StatusInternalServerError
		goto end
	}
	if len(ls.Links) == 0 {
		hr.HTML = safehtml.HTMLFromConstant("<div>No links for selection</div>")
		hr.HTTPStatus = http.StatusNoContent
		goto end
	}
	hr.HTML, err = linksetTemplate.ExecuteToHTML(htmlLinkset{
		apiURL:     params.Host.URL(),
		Links:      ls.Links,
		Label:      params.FilterLabel.String(),
		requestURI: params.RequestURI.String(),
		queryJSON:  params.FilterJSON(),
	})
	if err != nil {
		goto end
	}
end:
	return hr, err
}

func GetLinksetSuccessAlertHTML(ctx Context, linkIds []int64) HTMLResponse {
	linkIds = linkIds[:4]
	linkURLs, err := model.LoadLinkURLs(ctx, linkIds)
	if err != nil {
		slog.Error("Failed to get link URLs for %v", linkIds)
	}
	if len(linkURLs) > 3 {
		linkURLs = linkURLs[:4]
		linkURLs[3] = "..."
	}
	msg := "TODO: Set message here."
	alertHTML, _ := GetAlertHTML(ctx, AlertParams{
		OOB:  true,
		Type: SuccessAlert,
		Message: Message{
			Text:  msg,
			Items: linkURLs,
		},
	})
	return alertHTML
}
