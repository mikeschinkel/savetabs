package shared

import (
	"fmt"
)

func RootURL() string {
	return fmt.Sprintf("http://localhost:%d", DefaultPort)
}
