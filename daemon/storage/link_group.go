package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mikeschinkel/savetabs/daemon/shared"
	"github.com/mikeschinkel/savetabs/daemon/sqlc"
)

type LinkGroup struct {
	GroupName string `json:"group_name"`
	GroupSlug string `json:"group_slug"`
	GroupType string `json:"group_type"`
	LinkURL   string `json:"link_url"`
}

type MoveLinkToGroupArgs struct {
	LinkIds     []int64
	FromGroupId int64
	ToGroupId   int64
}

type LinkGroupsForJSON []LinkGroupForJSON
type LinkGroupForJSON struct {
	GroupRole string `json:"group_role"`
	GroupId   int64  `json:"group_id"`
	LinkId    int64  `json:"link_id"`
}

func (lgs LinkGroupsForJSON) JSON() string {
	b, err := json.Marshal(lgs)
	if err != nil {
		slog.Error("Convert to JSON for Link Group Move", "args", lgs, "error", err)
	}
	return string(b)
}

func newLinkGroupsForJSON(m MoveLinkToGroupArgs) LinkGroupsForJSON {
	lgs := make(LinkGroupsForJSON, len(m.LinkIds)*2)
	var n int
	for _, role := range []string{"from", "to"} {
		for _, linkId := range m.LinkIds {
			lgs[n] = LinkGroupForJSON{
				GroupRole: role,
				LinkId:    linkId,
			}
			switch role {
			case "from":
				lgs[n].GroupId = m.FromGroupId
			case "to":
				lgs[n].GroupId = m.ToGroupId
			}
			n++
		}
	}
	return lgs
}

func (m MoveLinkToGroupArgs) Error() error {
	return errors.Join(
		fmt.Errorf("from_group=%d, ", m.FromGroupId),
		fmt.Errorf("link_ids=%s, ", shared.Int64Slice(m.LinkIds).Join(",")),
		fmt.Errorf("to_group=%d, ", m.ToGroupId),
	)
}

type LinkGroupMoveStatus struct {
	sqlc.LoadLinkGroupMoveStatusRow
}

func (s LinkGroupMoveStatus) HasExceptions() (has bool) {
	if s.ErrorException {
		has = true
	}
	if s.ToGroupException {
		has = true
	}
	if s.FromGroupException {
		has = true
	}
	if s.FromLinkGroupExceptions != 0 {
		has = true
	}
	if s.ToLinkGroupExceptions != 0 {
		has = true
	}
	if s.LinkExceptions != 0 {
		has = true
	}
	return has
}

type LinkGroupMoveException struct {
	Exception     string
	FromGroupId   int64
	FromGroupName string
	LinkId        int64
	LinkURL       string
	ToGroupId     int64
	ToGroupName   string
}

type MoveLinksToGroupResult struct {
	Status     LinkGroupMoveStatus
	Exceptions []LinkGroupMoveException
	MovedIds   []int64
}

func MoveLinksToGroup(ctx Context, dbtx *NestedDBTX, args MoveLinkToGroupArgs) (result MoveLinksToGroupResult, err error) {
	err = execWithEnsuredNestedDBTX(dbtx, func(dbtx *NestedDBTX) (err error) {
		var exceptions []sqlc.ListLinkGroupMoveExceptionsRow
		me := shared.NewMultiErr()
		lgsj := newLinkGroupsForJSON(args)
		err = insertLinkGroupMoveFromVarJSON(ctx, dbtx, lgsj.JSON(), func(ctx Context, q *sqlc.Queries, ids LinkGroupMoveIds) (err error) {
			var status sqlc.LoadLinkGroupMoveStatusRow
			status, err = q.LoadLinkGroupMoveStatus(ctx, sqlc.LoadLinkGroupMoveStatusParams{
				FromIds: ids.FromIds,
				ToIds:   ids.ToIds,
			})
			if err != nil {
				goto end
			}
			result.Status = LinkGroupMoveStatus{status}
			if result.Status.HasExceptions() {
				exceptions, err = dbtx.DataStore.Queries(dbtx).ListLinkGroupMoveExceptions(ctx, sqlc.ListLinkGroupMoveExceptionsParams{
					FromIds: ids.FromIds,
					ToIds:   ids.ToIds,
				})
			}
			if status.ErrorException {
				// Errors exceptions mean we cannot continue
				goto end
			}
			// Some exceptions are only warnings and not errors
			//result.MovedIds, err = q.MoveLinksToGroup(ctx, sqlc.MoveLinksToGroupParams{
			//	GroupID:   args.ToGroupId,
			//	GroupID_2: args.FromGroupId,
			//	LinkIds:   args.LinkIds,
			//})
		end:
			return err
		})
		if err != nil {
			me.Add(err)
		}
		result.Exceptions = shared.ConvertSlice(exceptions, func(row sqlc.ListLinkGroupMoveExceptionsRow) LinkGroupMoveException {
			return LinkGroupMoveException{
				Exception:     row.Exception,
				FromGroupId:   row.FromGroupID,
				FromGroupName: row.FromGroupName,
				LinkId:        row.LinkID,
				LinkURL:       row.LinkUrl,
				ToGroupId:     row.ToGroupID,
				ToGroupName:   row.ToGroupName,
			}
		})
		return me.Err()
	})
	return result, err
}

type LinkGroupMoveIds struct {
	FromIds []int64
	ToIds   []int64
}

func insertLinkGroupMoveFromVarJSON(ctx Context, dbtx *NestedDBTX, j string, fn func(ctx Context, q *sqlc.Queries, ids LinkGroupMoveIds) error) (err error) {
	return execWithJSON(ctx, dbtx, j, func(ctx Context, q *sqlc.Queries, id int64) error {
		var rows []sqlc.InsertLinkGroupMoveFromVarJSONRow
		var ids LinkGroupMoveIds
		var idCnt int

		rows, err = q.InsertLinkGroupMoveFromVarJSON(ctx, id)
		if err != nil {
			goto end
		}
		if len(rows) == 0 {
			slog.Warn("Failed to return IDs for insert link group move", "var_id", id, "json", j)
			goto end
		}
		idCnt = len(rows) / 2
		ids = LinkGroupMoveIds{
			FromIds: make([]int64, 0, idCnt),
			ToIds:   make([]int64, 0, idCnt),
		}
		for _, row := range rows {
			if row.Role == "from" {
				ids.FromIds = append(ids.FromIds, row.ID)
			} else {
				ids.ToIds = append(ids.ToIds, row.ID)
			}
		}
		err = fn(ctx, q, ids)
		if err != nil {
			goto end
		}
		err = q.DeleteLinkGroupMove(ctx, append(ids.FromIds, ids.ToIds...))
		if err != nil {
			goto end
		}
	end:
		return err
	})
}
