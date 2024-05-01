package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"slices"
	"strings"

	"savetabs/augment"
	"savetabs/shared"
	"savetabs/sqlc"
)

// UpsertLinksWithGroups converts type for slice to type
func UpsertLinksWithGroups(ctx Context, db *sqlc.NestedDBTX, gs LinksWithGroupsGetSetter) error {
	return UpsertLinks(ctx, db, LinksWithGroups(gs.GetLinksWithGroups()))
}

type LinksWithGroups []LinkWithGroupGetSetter
type LinkWithGroup struct {
	LinkWithGroupGetSetter
}

func (links LinksWithGroups) Upsert(ctx context.Context, db *sqlc.NestedDBTX) error {
	var groupBytes []byte
	var metaBytes []byte
	var linkGroupBytes []byte
	var gg []group
	var rgs []linkGroup
	var mm []Meta
	var me = shared.NewMultiErr()

	slog.Info("Received from Chrome extension", "num_links", len(links))

	urls := links.urls()

	urlBytes, err := json.Marshal(urls)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "link"))
	}

	throttle()
	gg = links.groups()
	groupBytes, err = json.Marshal(gg)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "group"))
	}

	throttle()
	rgs = links.linkGroups()
	linkGroupBytes, err = json.Marshal(rgs)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "link_group"))
	}

	throttle()
	mm, err = links.metaFromURLs(urls)
	if err != nil {
		me.Add(err, ErrFailedToExtractMeta)
	}

	throttle()
	err = db.Exec(func(tx sqlc.DBTX) (err error) {
		// TODO: Need to use tx, somehow
		err = sqlc.UpsertLinks(ctx, db, string(urlBytes))
		if err != nil {
			me.Add(err, ErrFailedUpsertLinks)
		}

		throttle()
		err = sqlc.UpsertLinkGroups(ctx, db, string(linkGroupBytes))
		if err != nil {
			me.Add(err, ErrFailedUpsertLinkGroups)
		}

		throttle()
		err = sqlc.UpsertGroups(ctx, db, string(groupBytes))
		if err != nil {
			me.Add(err, ErrFailedUpsertGroups)
		}

		metaBytes, err = json.Marshal(mm)
		if err != nil {
			me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "meta"))
		}

		throttle()
		err = sqlc.UpsertMeta(ctx, db, string(metaBytes))
		if err != nil {
			me.Add(err, ErrFailedUpsertMeta)
		}
		slog.Info("Saved",
			"num_links", len(links),
			"num_link_groups", len(rgs),
			"num_groups", len(gg),
			"num_meta", len(mm),
		)
		return err
	})
	return me.Err()
}

func (links LinksWithGroups) urls() []string {
	var appended = make(map[string]struct{})
	var urls = make([]string, 0)
	for _, r := range links {
		if r.GetOriginalURL() == "" {
			continue
		}
		_, seen := appended[r.GetOriginalURL()]
		if seen {
			continue
		}
		appended[r.GetOriginalURL()] = struct{}{}
		urls = append(urls, r.GetOriginalURL())
	}
	return urls
}

func (links LinksWithGroups) groups() []group {
	var appended = make(map[int64]struct{})
	var groups = make([]group, 0)
	var keywords []string
	for _, g := range links {
		if g.GetGroup() == "" {
			continue
		}
		if g.GetGroup() == "<none>" {
			continue
		}
		if g.GetGroupType() == "" {
			continue
		}
		_, seen := appended[g.GetGroupId()]
		if seen {
			continue
		}
		gt := sqlc.GroupTypeFromName(g.GetGroupType())
		appended[g.GetGroupId()] = struct{}{}
		groups = append(groups, group{
			Id:   g.GetGroupId(),
			Type: gt,
			Name: g.GetGroup(),
			Slug: fmt.Sprintf("%s/%s", strings.ToLower(gt), shared.Slugify(g.GetGroup())),
		})
		keywords = augment.ParseKeywords(g.GetOriginalURL())
		groups = append(groups, shared.MapSliceFunc(keywords, func(kw string) group {
			return group{
				Type: "K",
				Name: kw,
				Slug: fmt.Sprintf("k/%s", shared.Slugify(kw)),
			}
		})...)
	}
	return groups
}

func (links LinksWithGroups) linkGroups() []linkGroup {
	var appended = make(map[int64]map[string]struct{})
	var rgs = make([]linkGroup, 0)
	for _, rg := range links {
		if rg.GetGroup() == "" {
			continue
		}
		if rg.GetGroupType() == "" {
			continue
		}
		_, seen := appended[rg.GetGroupId()]
		if !seen {
			appended[rg.GetGroupId()] = make(map[string]struct{})
		}
		if rg.GetOriginalURL() == "" {
			continue
		}
		_, seen = appended[rg.GetGroupId()][rg.GetOriginalURL()]
		if seen {
			continue
		}
		appended[rg.GetGroupId()][rg.GetOriginalURL()] = struct{}{}
		rgs = append(rgs, linkGroup{
			LinkURL:   rg.GetOriginalURL(),
			GroupName: rg.GetGroup(),
			GroupSlug: shared.Slugify(rg.GetGroup()),
			GroupType: sqlc.GroupTypeFromName(rg.GetGroupType()),
		})
	}
	return rgs
}

func appendMetaIfNotEmpty(kvs []Meta, u string, key, value string) []Meta {
	if key == "" {
		return kvs
	}
	if value == "" {
		return kvs
	}
	return append(kvs, Meta{
		Url:   u,
		Key:   key,
		Value: value,
	})
}

func (links LinksWithGroups) meta() (kvs []Meta, err error) {
	return links.metaFromURLs(links.urls())
}

func (links LinksWithGroups) metaFromURLs(urls []string) (kvs []Meta, err error) {
	var urlObj *url.URL
	kvs = make([]Meta, 0)
	for _, u := range urls {
		if u == "" {
			continue
		}
		urlObj, err = url.Parse(u)
		if err != nil {
			goto end
		}
		if !urlObj.IsAbs() {
			err = errors.Join(ErrUrlNotAbsolute, fmt.Errorf("url=%s", u))
		}
		kvs = appendMetaIfNotEmpty(kvs, u, "scheme", urlObj.Scheme)
		host := urlObj.Hostname()
		tld, sld, sub := extractDomains(host)
		kvs = append(kvs, []Meta{
			{
				Url:   u,
				Key:   "tld",
				Value: tld,
			},
			{
				Url:   u,
				Key:   "sld",
				Value: sld,
			},
			{
				Url:   u,
				Key:   "hostname",
				Value: host,
			},
		}...)
		kvs = appendMetaIfNotEmpty(kvs, u, "subdomain", sub)
		kvs = appendMetaIfNotEmpty(kvs, u, "path", urlObj.RawPath)
		kvs = appendMetaIfNotEmpty(kvs, u, "query", urlObj.RawQuery)
		kvs = appendMetaIfNotEmpty(kvs, u, "fragment", urlObj.RawFragment)
	}
end:
	return kvs, err
}

func extractDomains(host string) (tld, sld, sub string) {
	pos := strings.LastIndexByte(host, '.')
	if pos == -1 {
		goto end
	}
	if len(host) == 0 {
		goto end
	}
	tld = host[pos+1:]
	host = host[:pos]
	pos = strings.LastIndexByte(host, '.')
	if pos != -1 {
		sub = host[:pos]
		host = host[pos+1:]
	}
	sld = host + "." + tld

end:
	return tld, sld, sub
}

func sanitizeLinksWithGroups(links LinksWithGroups) (_ LinksWithGroups, err error) {
	var _err error
	for i := 0; i < len(links); i++ {
		var lwg LinkWithGroup
		lwg, _err = sanitizeLinkWithGroup(LinkWithGroup{links[i]})
		if _err == nil {
			links[i] = lwg.LinkWithGroupGetSetter
			continue
		}
		links = slices.Delete(links, i, i)
		i--
		_err = errors.Join(ErrFoundInLink, fmt.Errorf("index=%d", i), _err)
		if err == nil {
			err = _err
		} else {
			err = errors.Join(err, _err)
		}
	}
	return links, err
}

func sanitizeLinkWithGroup(link LinkWithGroup) (_ LinkWithGroup, err error) {
	var _link LinkGetSetter
	_link, err = sanitizeLink(link)
	if err != nil {
		goto end
	}
	link = _link.(LinkWithGroup)
	if link.GetGroup == nil || link.GetGroup() == "" {
		link.SetGroup("none")
	}
	if link.GetGroupType == nil || link.GetGroupType() == "" {
		link.SetGroupType("invalid")
	}
end:
	return link, err
}

func sanitizeLink(link LinkGetSetter) (_ LinkGetSetter, err error) {
	if link.GetOriginalURL() == "" {
		err = ErrUrlNotSpecified
		goto end
	}
	if link.GetId() == 0 {
		link.SetId(0)
	}
end:
	return link, err
}

func UpsertLinks(ctx Context, db *sqlc.NestedDBTX, links LinksWithGroups) error {
	links, err := sanitizeLinksWithGroups(links)
	if err != nil {
		goto end
	}
	err = db.Exec(func(dbtx sqlc.DBTX) error {
		return excludeUnwantedLinks(links).Upsert(ctx, db)
	})
end:
	return err
}

// excludeUnwantedLinks removes unwanted URLs such as "about:blank" and "chrome://*"
// TODO: Enhance to be end-user scriptable, ideally using two or more approaches:
//
//	https://github.com/expr-lang/expr | https://expr-lang.org/
//	https://github.com/google/cel-go
//	https://github.com/yuin/gopher-lua
//	https://github.com/dop251/goja
//	https://github.com/google/starlark-go
//	https://github.com/go-python/gpython
//	https://github.com/d5/tengo
//	https://github.com/mattn/anko
//	https://github.com/mikespook/goemphp
//	https://github.com/risor-io/risor
//	https://github.com/gentee/gentee
//	https://code.google.com/archive/p/gotcl/
//	https://github.com/krotik/ecal
//	https://github.com/elsaland/elsa
//	https://github.com/antonvolkoff/goluajit
//	https://github.com/risor-io/risor
func excludeUnwantedLinks(links LinksWithGroups) LinksWithGroups {
	wanted := make(LinksWithGroups, len(links))
	index := 0
	for _, link := range links {
		if link.GetOriginalURL == nil {
			continue
		}
		u := link.GetOriginalURL()
		switch {
		case u == "about:blank":
			continue
		case strings.HasPrefix(u, "chrome://"):
			continue
		}
		wanted[index] = link
		index++
	}
	return wanted
}
