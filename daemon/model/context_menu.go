package model

import (
	"savetabs/shared"
	"savetabs/storage"
)

func init() {
	shared.GroupContextMenuType.AddItem(shared.NewContextMenuItem("Rename"))
}

type ContextMenuArgs struct {
	ContextMenuType *shared.ContextMenuType
	DBId            int64
}

type ContextMenuItemArgs struct {
	ContextMenuArgs
	Name string
}

func UpdateContextMenuItemName(ctx Context, args ContextMenuItemArgs) (merged bool, err error) {
	switch args.ContextMenuType {
	case shared.GroupContextMenuType:
		merged, err = storage.UpdateGroupName(ctx, storage.GroupNameArgs{
			Name: args.Name,
			DBId: args.DBId,
		})
	}
	return merged, err
}
