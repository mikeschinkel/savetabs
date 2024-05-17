package ui

//import (
//	"fmt"
//	"strconv"
//	"strings"
//
//	"github.com/google/safehtml"
//	"savetabs/shared"
//)
//
//type htmlGroup struct {
//	Id        int64
//	Name      safehtml.HTML
//	Type      safehtml.Identifier
//	TypeName  safehtml.HTML
//	LinkCount int64
//	Links     []htmlLink
//	Host      string
//}
//
//func (g htmlGroup) LinksQueryParams() string {
//	return fmt.Sprintf("g=%s", g.Slug())
//}
//
//func (g htmlGroup) MenuItemType() safehtml.Identifier {
//	return g.Type
//}
//
//func (g htmlGroup) Slug() string {
//	return fmt.Sprintf("%s/%s",
//		strings.ToLower(g.Type.String()),
//		shared.Slugify(g.Name.String()),
//	)
//}
//
//func (g htmlGroup) URL() safehtml.URL {
//	return shared.MakeSafeURLf("http://%s/html/groups/%s",
//		g.Host,
//		g.Slug(),
//	)
//}
//
//func (g htmlGroup) Target() safehtml.Identifier {
//	return safehtml.IdentifierFromConstantPrefix(`group-links`,
//		strconv.FormatInt(g.Id, 10),
//	)
//}
//
//func (g htmlGroup) HTMLId() safehtml.Identifier {
//	return safehtml.IdentifierFromConstantPrefix(`group`,
//		strconv.FormatInt(g.Id, 10),
//	)
//}
