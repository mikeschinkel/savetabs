package ui

import (
	"github.com/google/safehtml"
	"github.com/mikeschinkel/savetabs/daemon/shared"
)

type IconState struct {
	value safehtml.Identifier
}

func newIconState(state string) IconState {
	return IconState{
		value: shared.MakeSafeId(state),
	}
}

//goland:noinspection GoUnusedGlobalVariable
var (
	ZeroStateIcon IconState
	BlankIcon     IconState = newIconState("blank")
	ExpandedIcon  IconState = newIconState("expanded")
	CollapsedIcon IconState = newIconState("collapsed")
)
