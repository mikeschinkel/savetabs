package ui

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/safehtml"
	"savetabs/shared"
)

var (
	draggableLink       = newDragDropTarget("link")
	droppableGroup      = newDragDropTarget("group")
	linkToGroupDragDrop = newDragDrop("link-to-group", draggableLink, droppableGroup)
)

type dragDrop struct {
	Name      safehtml.Identifier
	draggable *dragDropParticipant
	droppable *dragDropParticipant
}

func (d *dragDrop) String() string {
	return d.Name.String()
}

func newDragDrop(name string, draggable, droppable *dragDropParticipant) *dragDrop {
	return &dragDrop{
		Name:      shared.MakeSafeId(name),
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
	DragIds []int64
	DropId  int64
}

func newDropItem(dd *dragDrop, dropId int64) dragDropItem {
	return dragDropItem{
		dragDrop: dd,
		DragIds:  make([]int64, 0),
		DropId:   dropId,
	}
}

func newDragItem(dd *dragDrop, ids []int64) dragDropItem {
	return dragDropItem{
		dragDrop: dd,
		DragIds:  ids,
	}
}

func (d *dragDropItem) DragSources() string {
	return strings.Join(shared.ConvertSlice(d.DragIds, func(id int64) string {
		return fmt.Sprintf("%s:%d", d.draggable.Name, id)
	}), " ")
}

func (d *dragDropItem) DropTarget() string {
	return fmt.Sprintf("%s:%d", d.droppable.Name, d.DropId)
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

func dragDropTargetByName(name string) (_ *dragDropParticipant, err error) {
	mt, ok := dragDropTargetMap[strings.ToLower(name)]
	if !ok {
		err = errors.Join(ErrDragDropTargetNotFound, fmt.Errorf("target=%s", name))
	}
	return mt, err
}
