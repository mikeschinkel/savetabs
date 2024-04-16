package ui

type FilterType = string

var MetaFilter FilterType = "M"

const (
	GroupTypeFilter FilterType = "GT"

	BookmarkFilter FilterType = "B"
	TabGroupFilter FilterType = "G"
	TagFilter      FilterType = "T"
	CategoryFilter FilterType = "C"
	KeywordFilter  FilterType = "K"
	InvalidFilter  FilterType = "I"

	NoFilter FilterType = "_"
)

var FilterTypes = []FilterType{
	GroupTypeFilter,
	TagFilter,
	TabGroupFilter,
	CategoryFilter,
	KeywordFilter,
	BookmarkFilter,
	InvalidFilter,
}
