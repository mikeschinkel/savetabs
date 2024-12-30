package model

import (
	"net/url"

	"savetabs/shared"
	"savetabs/storage"
)

type AddLinksWithGroupsParams struct {
	Links []LinkWithGroup
}

type LinkGroup storage.LinkGroup

type LinksWithGroups struct {
	Links []LinkWithGroup
}

type LinkWithGroup struct {
	URL         *url.URL
	OriginalURL string
	Title       string
	GroupId     int64
	GroupSlug   string
	GroupType   shared.GroupType
	Group       string
}

func AddLinksWithGroupsIfNotExists(ctx Context, p AddLinksWithGroupsParams) (err error) {
	lswgs := LinksWithGroups{}
	lswgs.Links = p.Links

	links := shared.ConvertSlice(lswgs.UniqueLinks(), func(link LinkToAdd) storage.LinkToAdd {
		return storage.LinkToAdd{
			OriginalURL: link.OriginalURL,
			Title:       link.Title,
			Archived:    0,
			Deleted:     0,
		}
	})

	groups := shared.ConvertSlice(lswgs.Groups(), func(grp Group) storage.Group {
		return storage.Group{
			Id:       grp.Id,
			Name:     grp.Name,
			Type:     grp.Type.String(),
			Slug:     grp.Slug(),
			Archived: 0,
			Deleted:  0,
		}
	})

	groupedLinks := shared.ConvertSlice(lswgs.GroupedLinks(), func(grp LinkGroup) storage.LinkGroup {
		return storage.LinkGroup{
			GroupName: grp.GroupName,
			GroupSlug: grp.GroupSlug,
			GroupType: grp.GroupType,
			LinkURL:   grp.LinkURL,
		}
	})

	return storage.UpsertLinksWithGroups(ctx, storage.UpsertLinksWithGroupsParams{
		Links:        links,
		Groups:       groups,
		GroupedLinks: groupedLinks,
	})
}

// UniqueLinks collects all the unique links from the .Links property of the LinksWithGroups object.
func (lswgs LinksWithGroups) UniqueLinks() []LinkToAdd {
	var appended = make(map[string]struct{})
	var links = make([]LinkToAdd, len(lswgs.Links))
	var count = 0
	for i, link := range lswgs.Links {
		ou := link.OriginalURL
		if ou == "" {
			continue
		}
		_, seen := appended[ou]
		if seen {
			// De-dup URL from multiple tab groups.
			continue
		}
		count++
		appended[ou] = struct{}{}
		links[i] = LinkToAdd{
			OriginalURL: ou,
			Title:       link.Title,
			// TODO: MAYBE Parse the URL and add all other other info here
		}
	}
	return links[:count] // TODO: Verify this is not off by one.
}

func (lswgs LinksWithGroups) Groups() []Group {
	var appended = make(map[int64]struct{})
	var groups = make([]Group, 0)
	var keywords []string
	for _, g := range lswgs.Links {
		if g.Group == "" {
			continue
		}
		if g.Group == "<none>" {
			continue
		}
		if g.GroupType.Empty() {
			continue
		}
		_, seen := appended[g.GroupId]
		if seen {
			continue
		}
		appended[g.GroupId] = struct{}{}
		groups = append(groups, Group{
			Id:   g.GroupId,
			Type: g.GroupType,
			Name: g.Group,
		})
		keywords = ParseKeywords(g.OriginalURL)
		groups = append(groups, shared.ConvertSlice(keywords, func(kw string) Group {
			return Group{
				Type: shared.GroupTypeKeyword,
				Name: kw,
			}
		})...)
	}
	return groups
}

func (lswgs LinksWithGroups) GroupedLinks() []LinkGroup {
	var appended = make(map[int64]map[string]struct{})
	var lgs = make([]LinkGroup, 0)
	for _, lg := range lswgs.Links {
		if lg.Group == "" {
			continue
		}
		if lg.GroupType.Empty() {
			continue
		}
		_, seen := appended[lg.GroupId]
		if !seen {
			appended[lg.GroupId] = make(map[string]struct{})
		}
		if lg.OriginalURL == "" {
			continue
		}
		_, seen = appended[lg.GroupId][lg.OriginalURL]
		if seen {
			continue
		}
		appended[lg.GroupId][lg.OriginalURL] = struct{}{}
		lgs = append(lgs, LinkGroup{
			LinkURL:   lg.OriginalURL,
			GroupName: lg.Group,
			GroupSlug: lg.GroupSlug,
			GroupType: lg.GroupType.Upper(),
		})
	}
	return lgs
}
