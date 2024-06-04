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
	draggableLink          = newDragDropTarget("link")
	droppableMenuItem      = newDragDropTarget("menu-item")
	linkToMenuItemDragDrop = newDragDrop("link-to-menu-item", draggableLink, droppableMenuItem)
)

type dragDrop struct {
	Name      safehtml.Identifier
	draggable *dragDropParticipant
	droppable *dragDropParticipant
}

func (d *dragDrop) String() string {
	return d.Name.String()
}

func (d *dragDrop) DragSource() string {
	return d.draggable.Name
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
	FromDBIds []int64
	ToDBId    int64
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
