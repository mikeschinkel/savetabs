package shared

import (
	"fmt"
)

type ContextMenu struct {
	Id   int64
	Type *ContextMenuType
}

func (cm ContextMenu) Items() []ContextMenuItem {
	return cm.Type.Items
}

func (cm ContextMenu) String() string {
	return fmt.Sprintf("%s-%d", cm.Type.Name, cm.Id)
}

func NewContextMenu(cmt *ContextMenuType, id int64) *ContextMenu {
	return &ContextMenu{
		Id:   id,
		Type: cmt,
	}
}
