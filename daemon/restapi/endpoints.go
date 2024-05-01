package restapi

import (
	"fmt"
)

func LinksEndpoint() string {
	return fmt.Sprintf("%s/links", RootURL())
}

func LinkEndpoint(id int64) string {
	return fmt.Sprintf("%s/%d", LinksEndpoint(), id)
}
