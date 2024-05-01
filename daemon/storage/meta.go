package storage

type Meta struct {
	LinkId int64  `json:"link_id,omitempty"`
	Url    string `json:"url,omitempty"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
