package ui

import (
	"testing"

	"github.com/mikeschinkel/savetabs/daemon/model"
	"github.com/mikeschinkel/savetabs/daemon/shared"
)

var APIURL = shared.MakeSafeURL("http://localhost:8642")

type menuItemTest struct {
	name         string
	args         *HTMLMenuItemArgs
	mi           model.MenuItem
	want         HTMLMenuItem
	htmlId       string
	menuURL      string
	contentQuery string
}

type String struct {
	s string
}

func (s String) String() string {
	return s.s
}

func groupTypeKeywords() menuItemTest {
	gtm := NewHTMLMenu(HTMLMenuArgs{
		APIURL:   shared.MakeSafeURL("http://localhost:8642"),
		MenuType: shared.GroupTypeMenuType,
	})
	return menuItemTest{
		name: "Group Type: Keywords",
		mi: model.MenuItem{
			Menu:    model.NewMenu(shared.GroupTypeMenuType),
			LocalId: "k",
			Label:   "Keywords",
		},
		args: &HTMLMenuItemArgs{
			Parent: gtm,
		},
		want: HTMLMenuItem{
			parent:  gtm,
			localId: shared.GroupTypeKeyword.Lower(),
			Label:   shared.MakeSafeHTML("Keywords"),
		},
		htmlId:       "mi-gt-k",
		menuURL:      "gt--k",
		contentQuery: "gt=k",
	}
}
func groupKeywordNYTimes() menuItemTest {
	kwm := NewHTMLMenu(HTMLMenuArgs{
		APIURL:   shared.MakeSafeURL("http://localhost:8642"),
		MenuType: shared.KeywordMenuType,
	})
	nytMenu := shared.NewMenuType(shared.KeywordMenuType, String{"nytimes"}, nil)
	return menuItemTest{
		name: "Group Keyword: NYTimes",
		mi: model.MenuItem{
			Menu:    model.NewMenu(shared.KeywordMenuType),
			LocalId: "k",
			Label:   "Keywords",
		},
		args: &HTMLMenuItemArgs{
			Parent: kwm,
		},
		want: HTMLMenuItem{
			parent:   kwm,
			localId:  "nytimes",
			Label:    shared.MakeSafeHTML("New York Times"),
			menuType: nytMenu,
		},
		htmlId:       "mi-gt-k-nytimes",
		menuURL:      "gt--k/grp--nytimes",
		contentQuery: "gt=k&grp=nytimes",
	}
}

func Test_newMenuItem(t *testing.T) {
	tests := []menuItemTest{
		groupTypeKeywords(),
		groupKeywordNYTimes(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mi := newHTMLMenuItem(tt.mi, tt.args)
			if mi.HTMLId().String() != tt.htmlId {
				t.Errorf("HTMLId() = %v, want %v", mi.HTMLId(), tt.htmlId)
			}
			if mi.ChildMenuURL().String() != tt.menuURL {
				t.Errorf("ChildMenuURL() = %v, want %v", mi.ChildMenuURL(), tt.menuURL)
			}
			if mi.ContentQuery().String() != tt.contentQuery {
				t.Errorf("contentQuery() = %v, want %v", mi.ContentQuery(), tt.contentQuery)
			}
			//if got := newMenuItem(tt.args); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("newMenuItem() = %v, want %v", got, tt.want)
			//}
		})
	}
}

//
//func TestLoadMenuItems(t *testing.T) {
//	type args struct {
//		ctx Context
//		p   LoadMenuItemParams
//	}
//	tests := []struct {
//		name      string
//		args      args
//		wantItems MenuItems
//		wantErr   bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotItems, err := LoadMenuItems(tt.args.ctx, tt.args.p)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("LoadMenuItems() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(gotItems, tt.wantItems) {
//				t.Errorf("LoadMenuItems() gotItems = %v, want %v", gotItems, tt.wantItems)
//			}
//		})
//	}
//}
//
//func TestMenuItem_HTMLId(t *testing.T) {
//	type fields struct {
//		MenuItemParent shared.MenuItemParent
//		LocalId      string
//		Label        string
//		Type         shared.MenuType
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   safehtml.Identifier
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mi := MenuItem{
//				MenuItemParent: tt.fields.MenuItemParent,
//				LocalId:      tt.fields.LocalId,
//				Label:        tt.fields.Label,
//				Type:         tt.fields.Type,
//			}
//			if got := mi.HTMLId(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("HTMLId() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMenuItem_ItemURL(t *testing.T) {
//	type fields struct {
//		MenuItemParent shared.MenuItemParent
//		LocalId      string
//		Label        string
//		Type         shared.MenuType
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		wantU  safehtml.URL
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mi := MenuItem{
//				MenuItemParent: tt.fields.MenuItemParent,
//				LocalId:      tt.fields.LocalId,
//				Label:        tt.fields.Label,
//				Type:         tt.fields.Type,
//			}
//			if gotU := mi.ItemURL(); !reflect.DeepEqual(gotU, tt.wantU) {
//				t.Errorf("ItemURL() = %v, want %v", gotU, tt.wantU)
//			}
//		})
//	}
//}
//
//func TestMenuItem_Level(t *testing.T) {
//	type fields struct {
//		MenuItemParent shared.MenuItemParent
//		LocalId      string
//		Label        string
//		Type         shared.MenuType
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   int
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mi := MenuItem{
//				MenuItemParent: tt.fields.MenuItemParent,
//				LocalId:      tt.fields.LocalId,
//				Label:        tt.fields.Label,
//				Type:         tt.fields.Type,
//			}
//			if got := mi.Level(); got != tt.want {
//				t.Errorf("Level() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMenuItem_MenuType(t *testing.T) {
//	type fields struct {
//		MenuItemParent shared.MenuItemParent
//		LocalId      string
//		Label        string
//		Type         shared.MenuType
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   shared.MenuType
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mi := MenuItem{
//				MenuItemParent: tt.fields.MenuItemParent,
//				LocalId:      tt.fields.LocalId,
//				Label:        tt.fields.Label,
//				Type:         tt.fields.Type,
//			}
//			if got := mi.MenuType(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("MenuType() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMenuItem_Renew(t *testing.T) {
//	type fields struct {
//		MenuItemParent shared.MenuItemParent
//		LocalId      string
//		Label        string
//		Type         shared.MenuType
//	}
//	type args struct {
//		p MenuItemParams
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   MenuItem
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mi := MenuItem{
//				MenuItemParent: tt.fields.MenuItemParent,
//				LocalId:      tt.fields.LocalId,
//				Label:        tt.fields.Label,
//				Type:         tt.fields.Type,
//			}
//			if got := mi.Renew(tt.args.p); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Renew() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMenuItem_Root(t *testing.T) {
//	type fields struct {
//		MenuItemParent shared.MenuItemParent
//		LocalId      string
//		Label        string
//		Type         shared.MenuType
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   shared.MenuItemParent
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mi := MenuItem{
//				MenuItemParent: tt.fields.MenuItemParent,
//				LocalId:      tt.fields.LocalId,
//				Label:        tt.fields.Label,
//				Type:         tt.fields.Type,
//			}
//			if got := mi.Root(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Root() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMenuItem_ChildMenuURL(t *testing.T) {
//	type fields struct {
//		MenuItemParent shared.MenuItemParent
//		LocalId      string
//		Label        string
//		Type         shared.MenuType
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   safehtml.URL
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mi := MenuItem{
//				MenuItemParent: tt.fields.MenuItemParent,
//				LocalId:      tt.fields.LocalId,
//				Label:        tt.fields.Label,
//				Type:         tt.fields.Type,
//			}
//			if got := mi.ChildMenuURL(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("ChildMenuURL() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
