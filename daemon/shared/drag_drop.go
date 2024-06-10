package shared

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/safehtml"
)

var (
	DraggableLink       = NewDragDropTarget("link")
	DroppableGroup      = NewDragDropTarget("group")
	LinkToGroupDragDrop = NewDragDrop(DraggableLink, DroppableGroup)
)

type DragDrop struct {
	Name      safehtml.Identifier
	draggable *DragDropParticipant
	droppable *DragDropParticipant
}

func (d *DragDrop) String() string {
	return fmt.Sprintf("%s => %s", d.draggable, d.droppable)
}
func (d *DragDrop) Slug() string {
	return fmt.Sprintf("%s:%s", d.droppable.Name, d.draggable.Name)
}

var (
	dragDrops     = make([]*DragDrop, 0)
	dragDropMap   = make(map[string]*DragDrop)
	dragDropMutex sync.Mutex
)

func NewDragDrop(draggable, droppable *DragDropParticipant) *DragDrop {
	dd := &DragDrop{
		draggable: draggable,
		droppable: droppable,
	}
	dragDropMutex.Lock()
	defer dragDropMutex.Unlock()
	dragDrops = append(dragDrops, dd)
	dragDropMap[dd.Slug()] = dd
	return dd
}

type DragDropParticipant struct {
	Name  string
	Table string
}

type DragDropItem struct {
	*DragDrop
	DragParentId int64
	DragIds      []int64
	DropId       int64
}

func NewDropItem(dd *DragDrop, dropId int64) DragDropItem {
	return DragDropItem{
		DragDrop: dd,
		DragIds:  make([]int64, 0),
		DropId:   dropId,
	}
}

func NewDragItem(dd *DragDrop, dragParentId int64, ids []int64) DragDropItem {
	return DragDropItem{
		DragDrop:     dd,
		DragIds:      ids,
		DragParentId: dragParentId,
	}
}

func (d *DragDropItem) DragSources() string {
	return strings.Join(ConvertSlice(d.DragIds, func(id int64) string {
		return fmt.Sprintf("%s:%d", d.draggable.Name, id)
	}), " ")
}
func (d *DragDropItem) DragParent() string {
	return fmt.Sprintf("%s:%d", d.droppable.Name, d.DragParentId)
}

func (d *DragDropItem) DropTarget() string {
	return fmt.Sprintf("%s:%d", d.droppable.Name, d.DropId)
}

func (d *DragDropItem) DropTypes() string {
	// TODO: This is just one. Maybe in future we'll support many.
	return fmt.Sprintf("%s+%s", d.droppable.Name, d.draggable.Name)
}

var (
	dragDropTargets     = make([]*DragDropParticipant, 0)
	dragDropTargetMap   = make(map[string]*DragDropParticipant)
	dragDropTargetMutex sync.Mutex
)

func NewDragDropTarget(name string) *DragDropParticipant {
	ddt := &DragDropParticipant{
		Name:  name,
		Table: name,
	}
	dragDropTargetMutex.Lock()
	defer dragDropTargetMutex.Unlock()
	dragDropTargets = append(dragDropTargets, ddt)
	dragDropTargetMap[ddt.Name] = ddt
	return ddt
}

func DragDropByTypes(from, to Identifier) (dd *DragDrop, err error) {
	var ok bool

	ddType := fmt.Sprintf("%s:%s", to, from)
	dd, ok = dragDropMap[ddType]
	if !ok {
		err = errors.Join(ErrInvalidDragDropType, fmt.Errorf("type=%s", ddType))
	}
	return dd, err
}

//func dragDropTargetByName(name string) (_ *DragDropParticipant, err error) {
//	mt, ok := dragDropTargetMap[strings.ToLower(name)]
//	if !ok {
//		err = errors.Join(ErrDragDropTargetNotFound, fmt.Errorf("target=%s", name))
//	}
//	return mt, err
//}
