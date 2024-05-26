package shared

type FilterItem interface {
	HTMLId(MenuItemable) string
	ContentQuery(MenuItemable) string
	FilterType() *FilterType
	Label() string
	Filters() []any
	String() string
}
