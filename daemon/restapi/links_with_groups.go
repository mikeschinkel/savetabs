package restapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"savetabs/augment"
	"savetabs/shared"
	"savetabs/sqlc"
)

func (a *API) PostLinksWithGroups(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Find a better status result than "Bad Gateway"
		sendError(w, r, http.StatusBadGateway, err.Error())
		return
	}
	var data linksWithGroups
	err = json.Unmarshal(body, &data)
	if err != nil {
		sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	data, err = sanitizeLinksWithGroups(data)
	if err != nil {
		sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	ds := sqlc.GetDatastore()
	db, ok := ds.DB().(*sqlc.NestedDBTX)
	if !ok {
		sendError(w, r, http.StatusInternalServerError, "DB not a NestedDBTX")
		return
	}
	err = db.Exec(func(tx *sql.Tx) (err error) {
		err = upsertLinks(context.TODO(), ds, data)
		switch {
		case err == nil:
			goto end
		case errors.Is(err, ErrFailedToUnmarshal):
			sendError(w, r, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrFailedUpsertLinks):
			// TODO: Break out errors into different status for different reasons
			fallthrough
		default:
			sendError(w, r, http.StatusInternalServerError, err.Error())
		}
	end:
		return err
	})
}

type linksWithGroups LinksWithGroups

func (rr linksWithGroups) urls() []string {
	var appended = make(map[string]struct{})
	var urls = make([]string, 0)
	for _, r := range rr {
		if *r.Url == "" {
			continue
		}
		_, seen := appended[*r.Url]
		if seen {
			continue
		}
		appended[*r.Url] = struct{}{}
		urls = append(urls, *r.Url)
	}
	return urls
}

type group struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Slug string `json:"slug"`
}

func (rr linksWithGroups) groups() []group {
	var appended = make(map[int64]struct{})
	var groups = make([]group, 0)
	var keywords []string
	for _, g := range rr {
		if g.Group == nil {
			continue
		}
		if *g.Group == "" {
			continue
		}
		if *g.Group == "<none>" {
			continue
		}
		if g.GroupType == nil {
			continue
		}
		if *g.GroupType == "" {
			continue
		}
		_, seen := appended[*g.GroupId]
		if seen {
			continue
		}
		gt := groupTypeFromName(*g.GroupType)
		appended[*g.GroupId] = struct{}{}
		groups = append(groups, group{
			Id:   *g.GroupId,
			Type: gt,
			Name: *g.Group,
			Slug: fmt.Sprintf("%s/%s", strings.ToLower(gt), shared.Slugify(*g.Group)),
		})
		keywords = augment.ParseKeywords(*g.Url)
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

type linkGroup struct {
	GroupName string `json:"group_name"`
	GroupSlug string `json:"group_slug"`
	GroupType string `json:"group_type"`
	LinkURL   string `json:"link_url"`
}

func (rr linksWithGroups) linkGroups() []linkGroup {
	var appended = make(map[int64]map[string]struct{})
	var rgs = make([]linkGroup, 0)
	for _, rg := range rr {
		if rg.Group == nil {
			continue
		}
		if rg.GroupType == nil {
			continue
		}
		if *rg.Group == "" {
			continue
		}
		if *rg.GroupType == "" {
			continue
		}
		_, seen := appended[*rg.GroupId]
		if !seen {
			appended[*rg.GroupId] = make(map[string]struct{})
		}
		if rg.Url == nil {
			continue
		}
		if *rg.Url == "" {
			continue
		}
		_, seen = appended[*rg.GroupId][*rg.Url]
		if seen {
			continue
		}
		appended[*rg.GroupId][*rg.Url] = struct{}{}
		rgs = append(rgs, linkGroup{
			LinkURL:   *rg.Url,
			GroupName: *rg.Group,
			GroupSlug: shared.Slugify(*rg.Group),
			GroupType: groupTypeFromName(*rg.GroupType),
		})
	}
	return rgs
}

func groupTypeFromName(n string) (t string) {
	switch strings.ToLower(n) {
	case "category":
		t = "C"
	case "keyword":
		t = "K"
	case "tag":
		t = "T"
	case "tabgroup", "tab-group":
		t = "G"
	default:
		t = "I"
	}
	return t
}

func throttle() {
	time.Sleep(time.Second)
}

func upsertLinks(ctx context.Context, ds sqlc.DataStore, rr linksWithGroups) error {
	var groupBytes []byte
	var metadataBytes []byte
	var linkGroupBytes []byte
	var gg []group
	var rgs []linkGroup
	var mm []metadata
	var me = newMultiErr()

	slog.Info("Received from Chrome extension",
		"num_links", len(rr),
		"time", time.Now().Format(time.DateTime))

	urls := rr.urls()

	urlBytes, err := json.Marshal(urls)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "link"))
	}

	throttle()
	gg = rr.groups()
	groupBytes, err = json.Marshal(gg)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "group"))
	}

	throttle()
	rgs = rr.linkGroups()
	linkGroupBytes, err = json.Marshal(rgs)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "link_group"))
	}

	throttle()
	mm, err = rr.metadataFromURLs(urls)
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
	slog.Info("Received from Chrome extension",
		"num_links", len(rr),
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

func (rr linksWithGroups) metadata() (kvs []metadata, err error) {
	return rr.metadataFromURLs(rr.urls())
}

func (rr linksWithGroups) metadataFromURLs(urls []string) (kvs []metadata, err error) {
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

func sanitizeLinksWithGroups(data linksWithGroups) (_ linksWithGroups, err error) {
	for i := 0; i < len(data); i++ {
		rg := data[i]
		if rg.Url == nil || *rg.Url == "" {
			if err == nil {
				err = errors.Join(ErrUrlNotSpecified, fmt.Errorf("error found in link index %d", i))
			} else {
				err = errors.Join(err, fmt.Errorf("error found in link index %d", i))
			}
			data = slices.Delete(data, i, i)
			i--
			continue
		}
		if rg.Id == nil {
			data[i].Id = ptr[int64](0)
		}
		if rg.Group == nil || *rg.Group == "" {
			data[i].Group = ptr("none")
		}
		if rg.GroupType == nil || *rg.GroupType == "" {
			data[i].GroupType = ptr("invalid")
		}
	}
	return data, err
}
