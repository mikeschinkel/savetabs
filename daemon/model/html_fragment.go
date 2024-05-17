package model

type HTMLFragment struct {
	value string
}

func (hf HTMLFragment) String() string {
	return hf.value
}

func ParseHTMLFragment(h string) (_ HTMLFragment, err error) {
	// TODO: Add some validation of HTML
	return HTMLFragment{value: h}, err
}
