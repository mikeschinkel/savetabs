package shared

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	GroupContextMenuName     = "group"
	GroupTypeContextMenuName = "group_type"
)

var (
	GroupContextMenuType     = newContextMenuType(GroupContextMenuName)
	GroupTypeContextMenuType = newContextMenuType(GroupTypeContextMenuName)
)

var (
	contextMenuTypes     = make([]*ContextMenuType, 0)
	contextMenuTypeMap   = make(map[string]*ContextMenuType)
	contextMenuTypeMutex sync.Mutex
)

type ContextMenuType struct {
	Name  string
	table string
	Items []ContextMenuItem
}

func (cmt *ContextMenuType) AddItem(item ContextMenuItem) {
	cmt.Items = append(cmt.Items, item)
}

type ContextMenuItem struct {
	Label string
}

func NewContextMenuItem(label string) ContextMenuItem {
	return ContextMenuItem{Label: label}
}

func (cmt ContextMenuType) Table() string {
	return cmt.table
}

func newContextMenuType(name string) *ContextMenuType {
	cmt := &ContextMenuType{
		Name:  name,
		table: name,
	}
	contextMenuTypeMutex.Lock()
	defer contextMenuTypeMutex.Unlock()
	contextMenuTypes = append(contextMenuTypes, cmt)
	contextMenuTypeMap[cmt.Name] = cmt
	return cmt
}

func ContextMenuTypeByName(name string) (_ *ContextMenuType, err error) {
	mt, ok := contextMenuTypeMap[strings.ToLower(name)]
	if !ok {
		err = errors.Join(ErrContextMenuTypeNotFound, fmt.Errorf("name=%s", name))
	}
	return mt, err
}
