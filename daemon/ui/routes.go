package ui

// Routes is only here to provide easy access via IDE to each of our route's funcs.
//
//goland:noinspection GoUnusedGlobalVariable
var Routes = struct {
	GetBrowseHTML          func(Context, string) ([]byte, error)
	GetMenuHTML            func(Context, string) ([]byte, error)
	GetLinksHTML           func(Context, string, FilterValueGetter) ([]byte, error)
	GetGroupHTML           func(Context, string, string, string) ([]byte, error)
	GetGroupTypeGroupsHTML func(Context, string, string) ([]byte, error)
	GetMenuItemHTML        func(Context, string, string) ([]byte, error)
}{
	GetMenuHTML:            GetMenuHTML,
	GetLinksHTML:           GetLinksHTML,
	GetBrowseHTML:          GetBrowseHTML,
	GetGroupHTML:           GetGroupHTML,
	GetGroupTypeGroupsHTML: GetGroupTypeGroupsHTML,
	GetMenuItemHTML:        GetMenuItemHTML,
}
