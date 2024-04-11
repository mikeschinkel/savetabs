package restapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"savetabs/augment"
	"savetabs/shared"
	"savetabs/sqlc"
)

func (a *API) PostResourcesWithGroups(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Find a better status result than "Bad Gateway"
		sendError(w, r, http.StatusBadGateway, err.Error())
		return
	}
	var data resourcesWithGroups
	err = json.Unmarshal(body, &data)
	if err != nil {
		sendError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	data, err = sanitizeResourcesWithGroups(data)
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
		err = upsertResources(context.TODO(), ds, data)
		switch {
		case err == nil:
			goto end
		case errors.Is(err, ErrFailedToUnmarshal):
			sendError(w, r, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrFailedUpsertResources):
			// TODO: Break out errors into different status for different reasons
			fallthrough
		default:
			sendError(w, r, http.StatusInternalServerError, err.Error())
		}
	end:
		return err
	})
}

type resourcesWithGroups ResourcesWithGroups

func (rr resourcesWithGroups) urls() []string {
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

func (rr resourcesWithGroups) groups() []group {
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
		appended[*g.GroupId] = struct{}{}
		groups = append(groups, group{
			Id:   *g.GroupId,
			Type: groupTypeFromName(*g.GroupType),
			Name: *g.Group,
			Slug: shared.Slugify(*g.Group),
		})

		keywords = augment.ParseKeywords(*g.Url)
		groups = append(groups, shared.MapSliceFunc(keywords, func(kw string) group {
			return group{
				Type: "K",
				Name: kw,
				Slug: shared.Slugify(kw),
			}
		})...)
	}
	return groups
}

type resourceGroup struct {
	GroupName   string `json:"group_name"`
	GroupSlug   string `json:"group_slug"`
	GroupType   string `json:"group_type"`
	ResourceURL string `json:"resource_url"`
}

func (rr resourcesWithGroups) resourceGroups() []resourceGroup {
	var appended = make(map[int64]map[string]struct{})
	var rgs = make([]resourceGroup, 0)
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
		rgs = append(rgs, resourceGroup{
			ResourceURL: *rg.Url,
			GroupName:   *rg.Group,
			GroupSlug:   shared.Slugify(*rg.Group),
			GroupType:   groupTypeFromName(*rg.GroupType),
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

func upsertResources(ctx context.Context, ds sqlc.DataStore, rr resourcesWithGroups) error {
	var groupBytes []byte
	var keyValueBytes []byte
	var resourceGroupBytes []byte
	var gg []group
	var rgs []resourceGroup
	var kvs []keyValue
	var me = newMultiErr()

	log.Printf("Received new batch of resources from Chrome extension at %s",
		time.Now().Format(time.DateTime))

	urls := rr.urls()

	urlBytes, err := json.Marshal(urls)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "resource"))
	}

	gg = rr.groups()
	groupBytes, err = json.Marshal(gg)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "group"))
	}

	rgs = rr.resourceGroups()
	resourceGroupBytes, err = json.Marshal(rgs)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "resource_group"))
	}

	kvs, err = rr.keyValuesFromURLs(urls)
	if err != nil {
		me.Add(err, ErrFailedToExtractKeyValues)
	}

	err = sqlc.UpsertResources(ctx, ds, string(urlBytes))
	if err != nil {
		me.Add(err, ErrFailedUpsertResources)
	}

	err = sqlc.UpsertResourceGroups(ctx, ds, string(resourceGroupBytes))
	if err != nil {
		me.Add(err, ErrFailedUpsertResourceGroups)
	}

	err = sqlc.UpsertGroups(ctx, ds, string(groupBytes))
	if err != nil {
		me.Add(err, ErrFailedUpsertGroups)
	}

	keyValueBytes, err = json.Marshal(kvs)
	if err != nil {
		me.Add(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "key_value"))
	}
	err = sqlc.UpsertKeyValues(ctx, ds, string(keyValueBytes))
	if err != nil {
		me.Add(err, ErrFailedUpsertKeyValues)
	}
	log.Printf("Received %d resources, %d resource-groups, %d groups, and %d key/values from Chrome extension",
		len(rr), len(rgs), len(gg), len(kvs))
	return me.Err()
}

type keyValue struct {
	Url   *string `json:"url"`
	Key   string  `json:"key"`
	Value string  `json:"value"`
}

func appendKeyValueIfNotEmpty(kvs []keyValue, u *string, key, value string) []keyValue {
	if key == "" {
		return kvs
	}
	if value == "" {
		return kvs
	}
	return append(kvs, keyValue{
		Url:   u,
		Key:   key,
		Value: value,
	})
}

func (rr resourcesWithGroups) keyValues() (kvs []keyValue, err error) {
	return rr.keyValuesFromURLs(rr.urls())
}

func (rr resourcesWithGroups) keyValuesFromURLs(urls []string) (kvs []keyValue, err error) {
	var urlObj *url.URL
	kvs = make([]keyValue, 0)
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
		kvs = appendKeyValueIfNotEmpty(kvs, &u, "scheme", urlObj.Scheme)
		host := urlObj.Hostname()
		tld, sld, sub := extractDomains(host)
		kvs = append(kvs, []keyValue{
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
		kvs = appendKeyValueIfNotEmpty(kvs, &u, "subdomain", sub)
		kvs = appendKeyValueIfNotEmpty(kvs, &u, "path", urlObj.RawPath)
		kvs = appendKeyValueIfNotEmpty(kvs, &u, "query", urlObj.RawQuery)
		kvs = appendKeyValueIfNotEmpty(kvs, &u, "fragment", urlObj.RawFragment)
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

func sanitizeResourcesWithGroups(data resourcesWithGroups) (_ resourcesWithGroups, err error) {
	for i := 0; i < len(data); i++ {
		rg := data[i]
		if rg.Url == nil || *rg.Url == "" {
			if err == nil {
				err = errors.Join(ErrUrlNotSpecified, fmt.Errorf("error found in resource index %d", i))
			} else {
				err = errors.Join(err, fmt.Errorf("error found in resource index %d", i))
			}
			data = slices.Delete(data, i, i)
			i--
			continue
		}
		if rg.Id == nil {
			data[i].Id = ptr[int64](0)
		}
		if rg.Group == nil || *rg.Group == "" {
			data[i].Group = ptr("<none>")
		}
		if rg.GroupType == nil || *rg.GroupType == "" {
			data[i].GroupType = ptr("invalid")
		}
	}
	return data, err
}
