package storage

import (
	"context"
	"database/sql"
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
func UpsertLinksWithGroups(ctx Context, gs LinksWithGroupsGetSetter) error {
	return LinksWithGroups(gs.GetLinksWithGroups()).Upsert(ctx)
}

type LinksWithGroups []LinkWithGroupGetSetter

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

func (links LinksWithGroups) UpsertLinks(ctx context.Context, ds sqlc.DataStore) error {
	var groupBytes []byte
	var metadataBytes []byte
	var linkGroupBytes []byte
	var gg []group
	var rgs []linkGroup
	var mm []metadata
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
	mm, err = links.metadataFromURLs(urls)
	if err != nil {
		me.Add(err, ErrFailedToExtractMetadata)
	}

	throttle()
	err = sqlc.UpsertLinks(ctx, ds, string(urlBytes))
	if err != nil {
		me.Add(err, ErrFailedUpsertLinks)
	}

	throttle()
	err = sqlc.UpsertLinkGroups(ctx, ds, string(linkGroupBytes))
	if err != nil {
		me.Add(err, ErrFailedUpsertLinkGroups)
	}

	throttle()
	err = sqlc.UpsertGroups(ctx, ds, string(groupBytes))
	if err != nil {
		me.Add(err, ErrFailedUpsertGroups)
	}

	metadataBytes, err = json.Marshal(mm)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "metadata"))
	}

	throttle()
	err = sqlc.UpsertMetadata(ctx, ds, string(metadataBytes))
	if err != nil {
		me.Add(err, ErrFailedUpsertMetadata)
	}
	slog.Info("Saved",
		"num_links", len(links),
		"num_link_groups", len(rgs),
		"num_groups", len(gg),
		"num_meta", len(mm),
	)

	return me.Err()
}

type metadata struct {
	Url   *string `json:"url"`
	Key   string  `json:"key"`
	Value string  `json:"value"`
}

func appendMetadataIfNotEmpty(kvs []metadata, u *string, key, value string) []metadata {
	if key == "" {
		return kvs
	}
	if value == "" {
		return kvs
	}
	return append(kvs, metadata{
		Url:   u,
		Key:   key,
		Value: value,
	})
}

func (links LinksWithGroups) metadata() (kvs []metadata, err error) {
	return links.metadataFromURLs(links.urls())
}

func (links LinksWithGroups) metadataFromURLs(urls []string) (kvs []metadata, err error) {
	var urlObj *url.URL
	kvs = make([]metadata, 0)
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
		kvs = appendMetadataIfNotEmpty(kvs, &u, "scheme", urlObj.Scheme)
		host := urlObj.Hostname()
		tld, sld, sub := extractDomains(host)
		kvs = append(kvs, []metadata{
			{
				Url:   &u,
				Key:   "tld",
				Value: tld,
			},
			{
				Url:   &u,
				Key:   "sld",
				Value: sld,
			},
			{
				Url:   &u,
				Key:   "hostname",
				Value: host,
			},
		}...)
		kvs = appendMetadataIfNotEmpty(kvs, &u, "subdomain", sub)
		kvs = appendMetadataIfNotEmpty(kvs, &u, "path", urlObj.RawPath)
		kvs = appendMetadataIfNotEmpty(kvs, &u, "query", urlObj.RawQuery)
		kvs = appendMetadataIfNotEmpty(kvs, &u, "fragment", urlObj.RawFragment)
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

func sanitizeLinksWithGroups(data LinksWithGroups) (_ LinksWithGroups, err error) {
	for i := 0; i < len(data); i++ {
		rg := data[i]
		if rg.GetOriginalURL() == "" {
			if err == nil {
				err = errors.Join(ErrUrlNotSpecified, fmt.Errorf("error found in link index %d", i))
			} else {
				err = errors.Join(err, fmt.Errorf("error found in link index %d", i))
			}
			data = slices.Delete(data, i, i)
			i--
			continue
		}
		if rg.GetId() == 0 {
			data[i].SetId(0)
		}
		if rg.GetGroup == nil || rg.GetGroup() == "" {
			data[i].SetGroup("none")
		}
		if rg.GetGroupType == nil || rg.GetGroupType() == "" {
			data[i].SetGroupType("invalid")
		}
	}
	return data, err
}

func (links LinksWithGroups) Upsert(ctx Context) error {
	var ds sqlc.DataStore
	var db *sqlc.NestedDBTX
	var ok bool

	links, err := sanitizeLinksWithGroups(links)
	if err != nil {
		goto end
	}
	ds = sqlc.GetDatastore()
	db, ok = ds.DB().(*sqlc.NestedDBTX)
	if !ok {
		err = ErrDBNotANestedDCTX
		goto end
	}
	err = db.Exec(func(tx *sql.Tx) (err error) {
		// TODO: Need to use tx, somehow
		return excludeUnwantedLinks(links).UpsertLinks(ctx, ds)
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
