package ui

import (
	"fmt"
	"strings"
	"sync"

	"github.com/google/safehtml"
	"savetabs/shared"
)

var (
	draggableLink       = newDragDropTarget("link")
	droppableGroup      = newDragDropTarget("group")
	linkToGroupDragDrop = newDragDrop(draggableLink, droppableGroup)
)

type dragDrop struct {
	Name      safehtml.Identifier
	draggable *dragDropParticipant
	droppable *dragDropParticipant
}

func (d *dragDrop) String() string {
	return fmt.Sprintf("%s => %s", d.draggable, d.droppable)
}

func newDragDrop(draggable, droppable *dragDropParticipant) *dragDrop {
	return &dragDrop{
		draggable: draggable,
		droppable: droppable,
	}
}

type dragDropParticipant struct {
	Name  string
	Table string
}

type dragDropItem struct {
	*dragDrop
	DragParentId int64
	DragIds      []int64
	DropId       int64
}

func newDropItem(dd *dragDrop, dropId int64) dragDropItem {
	return dragDropItem{
		dragDrop: dd,
		DragIds:  make([]int64, 0),
		DropId:   dropId,
	}
}

func newDragItem(dd *dragDrop, dragParentId int64, ids []int64) dragDropItem {
	return dragDropItem{
		dragDrop:     dd,
		DragIds:      ids,
		DragParentId: dragParentId,
	}
}

func (d *dragDropItem) DragSources() string {
	return strings.Join(shared.ConvertSlice(d.DragIds, func(id int64) string {
		return fmt.Sprintf("%s:%d", d.draggable.Name, id)
	}), " ")
}
func (d *dragDropItem) DragParent() string {
	return fmt.Sprintf("%s:%d", d.droppable.Name, d.DragParentId)
}

func (d *dragDropItem) DropTarget() string {
	return fmt.Sprintf("%s:%d", d.droppable.Name, d.DropId)
}

func (d *dragDropItem) DropTypes() string {
	// TODO: This is just one. Maybe in future we'll support many.
	return fmt.Sprintf("%s+%s", d.droppable.Name, d.draggable.Name)
}

var (
	dragDropTargets     = make([]*dragDropParticipant, 0)
	dragDropTargetMap   = make(map[string]*dragDropParticipant)
	dragDropTargetMutex sync.Mutex
)

func newDragDropTarget(name string) *dragDropParticipant {
	ddt := &dragDropParticipant{
		Name:  name,
		Table: name,
	}
	dragDropTargetMutex.Lock()
	defer dragDropTargetMutex.Unlock()
	dragDropTargets = append(dragDropTargets, ddt)
	dragDropTargetMap[ddt.Name] = ddt
	return ddt
}

//func dragDropTargetByName(name string) (_ *dragDropParticipant, err error) {
//	mt, ok := dragDropTargetMap[strings.ToLower(name)]
//	if !ok {
//		err = errors.Join(ErrDragDropTargetNotFound, fmt.Errorf("target=%s", name))
//	}
//	return mt, err
//}
