package model

import (
	"savetabs/shared"
)

func init() {
	shared.GroupContextMenuType.AddItem(shared.NewContextMenuItem("Rename"))
	shared.GroupTypeContextMenuType.AddItem(shared.NewContextMenuItem("Rename"))
}
