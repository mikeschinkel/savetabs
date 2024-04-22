package storage

type linkGroup struct {
	GroupName string `json:"group_name"`
	GroupSlug string `json:"group_slug"`
	GroupType string `json:"group_type"`
	LinkURL   string `json:"link_url"`
}
