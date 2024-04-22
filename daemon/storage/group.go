package storage

type group struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Slug string `json:"slug"`
}
