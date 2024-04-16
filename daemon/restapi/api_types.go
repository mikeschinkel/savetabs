// Package restapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.1-0.20240325090356-a14414f04fdd DO NOT EDIT.
package restapi

// Error defines model for Error.
type Error struct {
	// Code Error code
	Code int32 `json:"code"`

	// Message Error message
	Message string `json:"message"`
}

// Id Primary key identifier
type Id = int64

// IdObject defines model for IdObject.
type IdObject struct {
	// Id Primary key identifier
	Id *Id `json:"id,omitempty"`
}

// IdObjects defines model for IdObjects.
type IdObjects = []IdObject

// Link defines model for Link.
type Link struct {
	Id    *int64  `json:"id,omitempty"`
	Title *string `json:"title,omitempty"`
	Url   *string `json:"url,omitempty"`
}

// LinkWithGroup defines model for LinkWithGroup.
type LinkWithGroup struct {
	Group     *string `json:"group,omitempty"`
	GroupId   *int64  `json:"groupId,omitempty"`
	GroupType *string `json:"groupType,omitempty"`
	Id        *int64  `json:"id,omitempty"`
	Title     *string `json:"title,omitempty"`
	Url       *string `json:"url,omitempty"`
}

// LinksWithGroups defines model for LinksWithGroups.
type LinksWithGroups = []LinkWithGroup

// BookmarkFilter defines model for BookmarkFilter.
type BookmarkFilter = []string

// CategoryFilter defines model for CategoryFilter.
type CategoryFilter = []string

// GroupSlug defines model for GroupSlug.
type GroupSlug = string

// GroupType defines model for GroupType.
type GroupType = string

// GroupTypeFilter defines model for GroupTypeFilter.
type GroupTypeFilter = []string

// GroupTypeName defines model for GroupTypeName.
type GroupTypeName = string

// KeywordFilter defines model for KeywordFilter.
type KeywordFilter = []string

// MenuItem defines model for MenuItem.
type MenuItem = string

// MetadataFilter defines model for MetadataFilter.
type MetadataFilter map[string]string

// TabGroupFilter defines model for TabGroupFilter.
type TabGroupFilter = []string

// TagFilter defines model for TagFilter.
type TagFilter = []string

// UnexpectedError defines model for UnexpectedError.
type UnexpectedError = Error

// GetLinksParams defines parameters for GetLinks.
type GetLinksParams struct {
	// Gt Links for a Group Type
	Gt *GroupTypeFilter `form:"gt,omitempty" json:"gt,omitempty"`

	// G TabGroup links by tags
	G *TabGroupFilter `form:"g,omitempty" json:"g,omitempty"`

	// C Category links by categories
	C *CategoryFilter `form:"c,omitempty" json:"c,omitempty"`

	// T Tag links by tags
	T *TagFilter `form:"t,omitempty" json:"t,omitempty"`

	// K Keyword filter for Links
	K *KeywordFilter `form:"k,omitempty" json:"k,omitempty"`

	// B Bookmark filter for Links
	B *BookmarkFilter `form:"b,omitempty" json:"b,omitempty"`

	// M Key/Value metadata filter for Links
	M *MetadataFilter `form:"m,omitempty" json:"m,omitempty"`
}

// PostLinksWithGroupsJSONRequestBody defines body for PostLinksWithGroups for application/json ContentType.
type PostLinksWithGroupsJSONRequestBody = LinksWithGroups