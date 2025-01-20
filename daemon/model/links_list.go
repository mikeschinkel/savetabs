package model

import (
	"github.com/mikeschinkel/savetabs/daemon/shared"
	"github.com/mikeschinkel/savetabs/daemon/storage"
)

type ListLinksParams = storage.ListLinksArgs

func ListLinks(ctx Context, p ListLinksParams) ([]Link, error) {
	links, err := storage.ListLinks(ctx, storage.ListLinksArgs(p))
	return shared.ConvertSlice(links, func(link storage.Link) Link {
		return Link(link)
	}), err
}
