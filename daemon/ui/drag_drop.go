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
	draggable *dragDropTarget
	droppable *dragDropTarget
}

func (d *dragDrop) String() string {
	return d.Name.String()
}

func newDragDrop(name string, draggable, droppable *dragDropTarget) *dragDrop {
	return &dragDrop{
		Name:      shared.MakeSafeId(name),
		draggable: draggable,
		droppable: droppable,
	}
}

type dragDropTarget struct {
	Name  string
	Table string
}

type dragDropItem struct {
	*dragDrop
	FromDBIds []int64
	ToDBId    int64
}

var (
	dragDropTargets     = make([]*dragDropTarget, 0)
	dragDropTargetMap   = make(map[string]*dragDropTarget)
	dragDropTargetMutex sync.Mutex
)

func newDragDropTarget(name string) *dragDropTarget {
	ddt := &dragDropTarget{
		Name:  name,
		Table: name,
	}
	dragDropTargetMutex.Lock()
	defer dragDropTargetMutex.Unlock()
	dragDropTargets = append(dragDropTargets, ddt)
	dragDropTargetMap[ddt.Name] = ddt
	return ddt
}

func dragDropTargetByName(name string) (_ *dragDropTarget, err error) {
	mt, ok := dragDropTargetMap[strings.ToLower(name)]
	if !ok {
		err = errors.Join(ErrDragDropTargetNotFound, fmt.Errorf("target=%s", name))
	}
	return mt, err
}
