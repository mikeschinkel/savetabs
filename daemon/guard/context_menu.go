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

func (args ContextMenuArgs) contextMenu() (cm *ui.ContextMenu, err error) {
	var cmt *shared.ContextMenuType

	cmt, err = shared.ContextMenuTypeByName(args.ContextMenuType)
	if err != nil {
		goto end
	}
	cm = ui.NewContextMenu(cmt, args.Host)
	cm.DBid = args.DBId
end:
	return cm, err
}

func GetContextMenuHTML(ctx Context, args ContextMenuArgs) (_ HTMLResponse, err error) {
	var cm *ui.ContextMenu

	hr := ui.NewHTMLResponse()

	cm, err = args.contextMenu()
	if err != nil {
		goto end
	}
	hr, err = ui.GetContextMenuHTML(ctx, ui.ContextMenuArgs{
		ContextMenu: cm,
		DBId:        args.DBId,
		Items: shared.ConvertSlice(cm.Type.Items, func(item shared.ContextMenuItem) ui.ContextMenuItem {
			return ui.ContextMenuItem{
				Label:       shared.MakeSafeHTML(item.Label),
				ContextMenu: cm,
			}
		}),
	})
end:
	return HTMLResponse{hr}, err
}

func GetContextMenuRenameFormHTML(ctx Context, args ContextMenuArgs) (_ HTMLResponse, err error) {
	var cm *ui.ContextMenu

	hr := ui.NewHTMLResponse()

	cm, err = args.contextMenu()
	if err != nil {
		goto end
	}

	hr, err = ui.GetContextMenuRenameFormHTML(ctx, ui.ContextMenuArgs{
		ContextMenu: cm,
		DBId:        args.DBId,
	})
end:
	return HTMLResponse{hr}, err
}
