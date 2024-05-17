package ui

type JSON struct {
	value string
}

func (j JSON) String() string {
	return j.value
}
