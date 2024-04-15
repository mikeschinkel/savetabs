package augment

import (
	"context"
	"encoding/json"
	"testing"

	"savetabs/shared"
	"savetabs/sqlc"
)

const DBFile = "../data/savetabs.db"

type linkGroup struct {
	GroupName string `json:"group_name"`
	GroupType string `json:"group_type"`
	LinkURL   string `json:"link_url"`
}
type group struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func TestKeywordsFromURL(t *testing.T) {
	ctx := context.Background()
	ds, err := sqlc.Initialize(ctx, DBFile)
	if err != nil {
		t.Fatal(err.Error())
	}
	q := ds.Queries()
	rr, err := q.ListLinks(ctx)

	var ggg []group
	var rgs []resourceGroup
	seen := make(map[string]struct{}, 0)
	for _, r := range rr {
		keywords := ParseKeywords(r.Url.String)
		gg := shared.MapSliceFunc(keywords, func(kw string) group {
			return group{
				Name: kw,
				Type: "K",
			}
		})
		for _, g := range gg {
			rgs = append(rgs, resourceGroup{
				GroupName:   g.Name,
				GroupType:   g.Type,
				ResourceURL: r.Url.String,
			})
		}
		for _, g := range gg {
			if _, ok := seen[g.Name]; ok {
				continue
			}
			seen[g.Name] = struct{}{}
			ggg = append(ggg, g)
		}
	}
	groupBytes, err := json.Marshal(ggg)
	if err != nil {
		t.Log(err.Error())
		return
	}
	err = sqlc.UpsertGroups(ctx, ds, string(groupBytes))
	if err != nil {
		t.Log(err.Error())
		return
	}
	rgBytes, err := json.Marshal(rgs)
	if err != nil {
		t.Log(err.Error())
		return
	}
	err = sqlc.UpsertResourceGroups(ctx, ds, string(rgBytes))
	if err != nil {
		t.Log(err.Error())
		return
	}
}
