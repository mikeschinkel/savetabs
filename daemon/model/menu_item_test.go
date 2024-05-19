package model

import (
	"testing"

	"savetabs/shared"
)

type menuItemTest struct {
	name    string
	args    MenuItemParams
	want    MenuItem
	htmlId  string
	menuURL string
	itemURL string
}

type String struct {
	s string
}

func (s String) String() string {
	return s.s
}

func groupTypeKeywords() menuItemTest {
	gtm := NewMenu(&shared.GroupTypeMenuType)
	return menuItemTest{
		name: "Group Type: Keywords",
		args: MenuItemParams{
			LocalId: "k",
			Label:   "Keywords",
			Menu:    gtm,
			Type:    &shared.KeywordMenuType,
		},
		want: MenuItem{
			MenuItemable: gtm,
			LocalId:      shared.GroupTypeKeyword.Lower(),
			Label:        "Keywords",
			Type:         &shared.KeywordMenuType,
		},
		htmlId:  "mi-gt-k",
		menuURL: "/gt/k",
		itemURL: "?gt=k",
	}
}
func groupKeywordNYTimes() menuItemTest {
	kwm := NewMenu(&shared.KeywordMenuType)
	nytMenu := shared.NewMenuType(&shared.KeywordMenuType, String{"nytimes"}, nil)
	return menuItemTest{
		name: "Group Keyword: NYTimes",
		args: MenuItemParams{
			LocalId: "nytimes",
			Label:   "New York Times",
			Menu:    kwm,
			Type:    &nytMenu,
		},
		want: MenuItem{
			MenuItemable: kwm,
			LocalId:      shared.GroupTypeKeyword.Lower(),
			Label:        "New York Times",
			Type:         &nytMenu,
		},
		htmlId:  "mi-gt-k-nytimes",
		menuURL: "/gt/k/nytimes",
		itemURL: "?gt=k&grp=nytimes",
	}
}

func Test_newMenuItem(t *testing.T) {
	tests := []menuItemTest{
		groupTypeKeywords(),
		groupKeywordNYTimes(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mi := newMenuItem(tt.args)
			if mi.HTMLId().String() != tt.htmlId {
				t.Errorf("HTMLId() = %v, want %v", mi.HTMLId(), tt.htmlId)
			}
			if mi.SubmenuURL().String() != tt.menuURL {
				t.Errorf("SubmenuURL() = %v, want %v", mi.SubmenuURL(), tt.menuURL)
			}
			if mi.ItemURL().String() != tt.itemURL {
				t.Errorf("ItemURL() = %v, want %v", mi.ItemURL(), tt.itemURL)
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
//		MenuItemable shared.MenuItemable
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
//				MenuItemable: tt.fields.MenuItemable,
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
//		MenuItemable shared.MenuItemable
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
//				MenuItemable: tt.fields.MenuItemable,
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
//		MenuItemable shared.MenuItemable
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
//				MenuItemable: tt.fields.MenuItemable,
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
//		MenuItemable shared.MenuItemable
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
//				MenuItemable: tt.fields.MenuItemable,
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
//		MenuItemable shared.MenuItemable
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
//				MenuItemable: tt.fields.MenuItemable,
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
//		MenuItemable shared.MenuItemable
//		LocalId      string
//		Label        string
//		Type         shared.MenuType
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   shared.MenuItemable
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mi := MenuItem{
//				MenuItemable: tt.fields.MenuItemable,
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
//func TestMenuItem_SubmenuURL(t *testing.T) {
//	type fields struct {
//		MenuItemable shared.MenuItemable
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
//				MenuItemable: tt.fields.MenuItemable,
//				LocalId:      tt.fields.LocalId,
//				Label:        tt.fields.Label,
//				Type:         tt.fields.Type,
//			}
//			if got := mi.SubmenuURL(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("SubmenuURL() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
