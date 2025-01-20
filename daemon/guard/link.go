package guard

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/mikeschinkel/savetabs/daemon/model"
	"github.com/mikeschinkel/savetabs/daemon/shared"
)

type UpsertLinkArgs struct {
	Id    int64
	URL   string
	Title string
	HTML  string
}

func (link UpsertLinkArgs) Sanitize() (_ UpsertLinkArgs, err error) {
	return link, err
}

func (link UpsertLinkArgs) Validate() (err error) {
	if link.URL == "" {
		err = ErrUrlNotSpecified
		goto end
	}
end:
	return err
}

func UpsertLink(ctx Context, link UpsertLinkArgs) (linkId int64, err error) {
	var content model.Content
	var lu *url.URL

	err = link.Validate()
	if err != nil {
		goto end
	}

	link, err = link.Sanitize()
	if err != nil {
		goto end
	}

	content, err = ParseContent(link.HTML)
	if err != nil {
		goto end
	}

	lu, err = url.Parse(link.URL)
	if err != nil {
		goto end
	}

	linkId, err = model.AddLink(ctx, model.LinkToAdd{
		URL:         *lu,
		OriginalURL: link.URL,
		Title:       link.Title,
		Content:     content,
	})

end:
	return linkId, err
}

type LinkToLoad model.Link

func LoadLink(ctx Context, linkId int64) (LinkToLoad, error) {
	link, err := model.LoadLink(ctx, linkId)
	return LinkToLoad(link), err
}

//func LoadLinkURLs(ctx Context, linkIds []int64) (urls []string, err error) {
//	return model.LoadLinkURLs(ctx, linkIds)
//}

// ParseLinkIds accepts a slice of Link Ids as strings and converts to a slice of
// int64, returning an error containing a joined error for each failing id.
func ParseLinkIds(linkIds []string) (ids []int64, _ error) {
	var err error
	me := shared.NewMultiErr()
	ids = make([]int64, len(linkIds))
	for i, id := range linkIds {
		ids[i], err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			me.Add(err, fmt.Errorf("link_id=%s", id))
		}
	}
	return ids, me.Err()
}

type ListLinksParams = model.ListLinksParams

func ListLinks(ctx Context, p ListLinksParams) ([]LinkToLoad, error) {
	links, err := model.ListLinks(ctx, model.ListLinksParams(p))
	return shared.ConvertSlice(links, func(link model.Link) LinkToLoad {
		return LinkToLoad(link)
	}), err
}
