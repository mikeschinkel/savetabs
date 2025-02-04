package ui

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/safehtml"
	"github.com/mikeschinkel/savetabs/daemon/model"
	"github.com/mikeschinkel/savetabs/daemon/shared"
)

type htmlLinkset struct {
	apiURL      safehtml.URL
	Links       []htmlLink
	requestURI  safehtml.URL
	filterQuery shared.FilterQuery
}

func (ls htmlLinkset) HeaderHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`links-row-head`)
}
func (ls htmlLinkset) Label() safehtml.HTML {
	return shared.MakeSafeHTML(ls.filterQuery.Label())
}
func (ls htmlLinkset) FooterHTMLId() safehtml.Identifier {
	return safehtml.IdentifierFromConstant(`links-row-foot`)
}
func (ls htmlLinkset) URLQuery() safehtml.URL {
	parts := strings.Split(ls.requestURI.String()+"?", "?")
	return safehtml.URLSanitized("?" + parts[1])
}
func (ls htmlLinkset) FilterQuery() safehtml.URL {
	return shared.MakeSafeURL(ls.filterQuery.String())
}
func (ls htmlLinkset) NumLinks() int {
	return len(ls.Links)
}
func (ls htmlLinkset) HTMLLinksURL() string {
	return fmt.Sprintf("%s/html/linkset", ls.apiURL)
}
func (ls htmlLinkset) TableHeaderFooterHTML() safehtml.HTML {
	return safehtml.HTMLFromConstant(`
<th class="w-1"></th>
<th class="p-0.5">#</th>
<th class="p-0.5">
	<label>
		<input type="checkbox" @change="maybeConfirmCheckAll" class="link-check-all"> 
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

type LinksetArgs struct {
	shared.FilterQuery
	RequestURI safehtml.URL
	APIURL     safehtml.URL
}

func (lp LinksetArgs) FilterJSON() (j safehtml.JSON) {
	b, err := json.Marshal(lp.FilterQuery)
	if err != nil {
		goto end
	}
	j, err = shared.MakeSafeJSON(string(b))
end:
	if err != nil {
		j = shared.MakeEmptyObjectJSON()
		slog.Error("Failed to marshal FilterQuery",
			"filter_query", lp.FilterQuery,
			"err", err.Error(),
		)
	}
	return j
}

func GetLinksetHTML(ctx Context, args LinksetArgs) (hr HTMLResponse, err error) {
	var ls model.LinksetToLoad
	var rowNum int
	var htmlLS htmlLinkset

	hr = NewHTMLResponse()

	ls, err = model.LoadLinkset(ctx, model.LinksetToLoadParams(model.LinksetToLoadParams{
		FilterQuery: args.FilterQuery,
	}))
	if err != nil {
		hr.StatusCode = http.StatusInternalServerError
		goto end
	}
	if len(ls.Links) == 0 {
		// TODO: Change to using a dismissible error
		hr.HTML = safehtml.HTMLFromConstant("<div>No links for selection</div>")
		hr.StatusCode = http.StatusNoContent
		goto end
	}

	htmlLS = htmlLinkset{
		apiURL:      args.APIURL,
		requestURI:  args.RequestURI,
		filterQuery: args.FilterQuery,
	}
	htmlLS.Links = shared.ConvertSlice(ls.Links, func(link model.Link) htmlLink {
		rowNum++
		parentId := args.FilterQuery.ParentDBId
		return newHTMLLink(htmlLinkArgs{
			Link:     link,
			RowId:    rowNum,
			DragItem: shared.NewDragItem(shared.LinkToGroupDragDrop, parentId, []int64{link.Id}),
		})
	})

	hr.HTML, err = linksetTemplate.ExecuteToHTML(htmlLS)
	if err != nil {
		goto end
	}
end:
	return hr, err
}

func GetLinksetSuccessAlertHTML(ctx Context, linkIds []int64) HTMLResponse {
	if len(linkIds) > 3 {
		linkIds = linkIds[:4]
	}
	linkURLs, err := model.LoadLinkURLs(ctx, linkIds)
	if err != nil {
		slog.Error("Failed to get link URLs for %v", "link_ids", linkIds)
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
