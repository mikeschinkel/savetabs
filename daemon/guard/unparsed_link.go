package guard

import (
	"github.com/mikeschinkel/savetabs/daemon/model"
	"github.com/mikeschinkel/savetabs/daemon/shared"
)

type UnparsedLink model.UnparsedLink

func LatestUnparsedLinks(ctx Context) (links []UnparsedLink, err error) {
	var modelLinks []model.UnparsedLink

	modelLinks, err = model.LatestUnparsedLinks(ctx)
	if err != nil {
		goto end
	}
	links = shared.ConvertSlice(modelLinks, func(link model.UnparsedLink) UnparsedLink {
		return UnparsedLink(link)
	})
end:
	return links, err
}

func UpdateUnparsedLink(ctx Context, link UnparsedLink) error {
	return model.UpdateUnparsedLink(ctx, model.UnparsedLink(link))
}
