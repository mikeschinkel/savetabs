package model

import (
	"savetabs/shared"
	"savetabs/storage"
)

type ListLinksParams = storage.ListLinksArgs

func ListLinks(ctx Context, p ListLinksParams) ([]Link, error) {
	links, err := storage.ListLinks(ctx, storage.ListLinksArgs(p))
	return shared.ConvertSlice(links, func(link storage.Link) Link {
		return Link(link)
	}), err
}
