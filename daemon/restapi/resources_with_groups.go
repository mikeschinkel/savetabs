package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"savetabs/sqlc"
)

type resourcesWithGroups ResourcesWithGroups

func (rr resourcesWithGroups) urls() []string {
	var appended = make(map[string]struct{})
	var urls = make([]string, 0)
	for _, r := range rr {
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
}

func (rr resourcesWithGroups) groups() []group {
	var appended = make(map[int64]struct{})
	var groups = make([]group, 0)
	for _, g := range rr {
		if g.Group == nil {
			continue
		}
		if g.GroupType == nil {
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
		})
	}
	return groups
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
	var gg []group
	var kvs []keyValue

	log.Printf("Received new batch of resources from Chrome extension at %s",
		time.Now().Format(time.DateTime))

	urls := rr.urls()

	urlBytes, err := json.Marshal(urls)
	if err != nil {
		err = errors.Join(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "resource"))
		goto end
	}

	gg = rr.groups()
	groupBytes, err = json.Marshal(gg)
	if err != nil {
		err = errors.Join(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "group"))
		goto end
	}

	kvs, err = rr.keyValuesFromURLs(urls)
	if err != nil {
		err = errors.Join(err, ErrFailedToExtractKeyValues)
		goto end
	}

	log.Printf("Received %d resources from Chrome extension", len(rr))
	err = sqlc.UpsertResources(ctx, ds, string(urlBytes))
	if err != nil {
		err = errors.Join(err, ErrFailedUpsertResources)
		goto end
	}

	log.Printf("Received %d groups from Chrome extension", len(gg))
	err = sqlc.UpsertGroups(ctx, ds, string(groupBytes))
	if err != nil {
		err = errors.Join(err, ErrFailedUpsertGroups)
		goto end
	}

	log.Printf("Derived %d key/values from resources", len(kvs))

	for {
		//keyValueBytes, err = json.Marshal(kvs[:100])
		keyValueBytes, err = json.Marshal(kvs)
		//		kvs = kvs[100:]
		if err != nil {
			err = errors.Join(err, ErrFailedToUnmarshal, fmt.Errorf("table=%s", "key_value"))
			goto end
		}
		err = sqlc.UpsertKeyValues(ctx, ds, string(keyValueBytes))
		if err != nil {
			err = errors.Join(err, ErrFailedUpsertKeyValues)
			goto end
		}
		goto end
	}
end:
	return err
}

type keyValue struct {
	Url   *string `json:"url"`
	Key   string  `json:"key"`
	Value string  `json:"value"`
}

func appendKeyValueIfNotEmpty(kvs []keyValue, u *string, key, value string) []keyValue {
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
