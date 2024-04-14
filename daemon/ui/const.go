package ui

type GroupType string

var MetaType = "M"

const (
	BookmarkGroupType = "B"
	TabGroupGroupType = "G"
	TagGroupType      = "T"
	CategoryGroupType = "C"
	KeywordGroupType  = "K"
	InvalidGroupType  = "I"
)

var GroupTypes = []GroupType{
	TagGroupType,
	TabGroupGroupType,
	CategoryGroupType,
	KeywordGroupType,
	BookmarkGroupType,
	InvalidGroupType,
}
