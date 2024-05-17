package shared

import (
	"strings"
)

var actions = make(map[string]struct{})

type Action struct {
	value string
}

func (a Action) String() string {
	return a.value
}

func ValidateAction(value string) bool {
	_, ok := actions[value]
	return ok
}

func NewAction(value string) Action {
	return Action{
		value: strings.ToLower(value),
	}
}

func registerAction(value string) Action {
	actions[value] = struct{}{}
	return NewAction(value)
}

var (
	ArchiveAction = registerAction("archive")
	DeleteAction  = registerAction("delete")
)
