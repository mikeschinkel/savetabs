package restapi

import (
	"fmt"
)

// LinksEndpoint returns `http://{root_path}/links`
func LinksEndpoint() string {
	return fmt.Sprintf("%s/links", RootURL())
}

// LinkEndpoint returns `http://{root_path}/links/{link_id}`
func LinkEndpoint(id int64) string {
	return fmt.Sprintf("%s/%d", LinksEndpoint(), id)
}
