package guard

import (
	"fmt"
	"net/url"
	"strconv"

	"savetabs/model"
	"savetabs/shared"
)

type UpsertLink struct {
	Id    int64
	URL   string
	Title string
	HTML  string
}

func (link UpsertLink) Sanitize() (_ UpsertLink, err error) {
	return link, err
}

func (link UpsertLink) Validate() (err error) {
	if link.URL == "" {
		err = ErrUrlNotSpecified
		goto end
	}
end:
	return err
}

func LinkUpsert(ctx Context, link UpsertLink) (linkId int64, err error) {
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

type LoadLinksParams = model.LoadLinksParams

func LoadLinks(ctx Context, p LoadLinksParams) ([]LinkToLoad, error) {
	links, err := model.LoadLinks(ctx, model.LoadLinksParams(p))
	return shared.ConvertSlice(links, func(link model.Link) LinkToLoad {
		return LinkToLoad(link)
	}), err
}
