package model

import (
	"savetabs/shared"
	"savetabs/storage"
)

type LoadLinksParams = storage.LoadLinksParams

func LoadLinks(ctx Context, p LoadLinksParams) ([]Link, error) {
	links, err := storage.LoadLinks(ctx, storage.LoadLinksParams(p))
	return shared.ConvertSlice(links, func(link storage.Link) Link {
		return Link(link)
	}), err
}
