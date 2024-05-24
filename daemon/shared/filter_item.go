package shared

type FilterItem interface {
	FilterType() *FilterType
	Label() string
	Filters() []any
}
