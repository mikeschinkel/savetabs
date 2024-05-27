package guard

import (
	"savetabs/shared"
	"savetabs/ui"
)

type ContextMenuArgs struct {
	Host            string
	ContextMenuType string
	DBId            int64
}

func GetContextMenuHTML(ctx Context, args ContextMenuArgs) (_ HTMLResponse, err error) {
	var cm ui.ContextMenu

	hr := ui.NewHTMLResponse()
	cmt, err := shared.ContextMenuTypeByName(args.ContextMenuType)
	if err != nil {
		goto end
	}
	cm = ui.NewContextMenu(cmt)

	hr, err = ui.GetContextMenuHTML(ctx, ui.ContextMenuArgs{
		APIURL:      shared.MakeSafeAPIURL(args.Host),
		ContextMenu: cm,
		Items: shared.ConvertSlice(cmt.Items, func(item shared.ContextMenuItem) ui.ContextMenuItem {
			return ui.ContextMenuItem{
				Label:           shared.MakeSafeHTML(item.Label),
				ContextMenuType: cmt,
			}
		}),
	})
end:
	return HTMLResponse{hr}, err
}
