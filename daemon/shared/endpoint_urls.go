package shared

import (
	"fmt"
)

// LinksEndpointURL returns `http://{root_path}/links`
func LinksEndpointURL() string {
	return fmt.Sprintf("%s/links", RootURL())
}

// LinkEndpointURL returns `http://{root_path}/links/{link_id}`
func LinkEndpointURL(id int64) string {
	return fmt.Sprintf("%s/%d", LinksEndpointURL(), id)
}

// GroupsHTMLEndpointURL returns `http://{root_path}/html/groups`
func GroupsHTMLEndpointURL() string {
	return fmt.Sprintf("%s/html/groups", RootURL())
}

// GroupHTMLEndpointURL returns `http://{root_path}/html/groups/{group_id}`
//
//	{group_id} => lower(type) + "/" + slugify(name)
func GroupHTMLEndpointURL(id string) string {
	return fmt.Sprintf("%s/%s", GroupsHTMLEndpointURL(), id)
}
