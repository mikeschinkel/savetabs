package shared

type Label struct {
	value string
}

func (l Label) String() string {
	return l.value
}

func NewLabel(value string) Label {
	return Label{value: value}
}
