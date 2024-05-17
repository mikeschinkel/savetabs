package shared

type Host struct {
	value string
}

func (h Host) String() string {
	return h.value
}

func (h Host) URL() string {
	// TODO: Does this need to be more robust?
	return "http://" + h.value
}

func NewHost(value string) Host {
	return Host{value: value}
}
