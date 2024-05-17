package guard

type JSON struct {
	value string
}

func (j JSON) String() string {
	return j.value
}

func NewJSON(j string) JSON {
	return JSON{value: j}
}
