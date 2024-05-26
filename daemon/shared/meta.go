package shared

import (
	"fmt"
)

type Meta struct {
	Key   string
	Value string
}

func (m Meta) String() string {
	return fmt.Sprintf("%s=%s", m.Key, m.Value)
}
