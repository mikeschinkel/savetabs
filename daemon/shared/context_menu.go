package shared

type ContextMenu struct {
	Id   int64
	Type *ContextMenuType
}

func (cm ContextMenu) Items() []ContextMenuItem {
	return cm.Type.Items
}

func NewContextMenu(cmt *ContextMenuType, id int64) ContextMenu {
	return ContextMenu{
		Id:   id,
		Type: cmt,
	}
}
