package model

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"savetabs/shared"
	"savetabs/storage"
)

type Group struct {
	Id   int64
	Name string
	Type shared.GroupType
}

func (grp Group) Slug() string {
	return fmt.Sprintf("%s:%s", grp.Type.Lower(), shared.Slugify(grp.Name))
}

type Groups struct {
	GroupsArgs
	Groups []Group
}

type GroupsArgs storage.GroupsArgs

func NewGroups(groups storage.Groups) Groups {
	gs := make([]Group, len(groups.Groups))
	for i, grp := range groups.Groups {
		gt, err := shared.ParseGroupTypeByLetter(grp.Type)
		if err != nil {
			// Panic because upstream should have cause this, so that needs to be where it is
			// fixed, not here. Hence failing here is a programming error.
			panic(err.Error())
		}
		gs[i] = Group{
			Id:   grp.Id,
			Name: grp.Name,
			Type: gt,
		}
	}
	return Groups{
		GroupsArgs: GroupsArgs(groups.Args),
		Groups:     gs,
	}
}

func LoadGroupName(ctx Context, groupId int64) (name string, err error) {
	return storage.LoadGroupName(ctx, nil, groupId)
}

func LoadGroupIdBySlug(ctx Context, slug string) (int64, error) {
	return storage.LoadGroupIdBySlug(ctx, nil, slug)
}

func LoadGroups(ctx Context, params GroupsArgs) (groups Groups, err error) {
	var gs storage.Groups
	gs, err = storage.LoadGroups(ctx, storage.GroupsArgs(params))
	if err != nil {
		goto end
	}
end:
	return NewGroups(gs), err
}

type MoveLinkToGroupArgs storage.MoveLinkToGroupArgs

type LinkGroupMoveStatus storage.LinkGroupMoveStatus
type LinkGroupMoveExceptions []LinkGroupMoveException
type LinkGroupMoveException storage.LinkGroupMoveException

func (e LinkGroupMoveException) Error() (s string) {
	switch e.Exception {
	case shared.FromGroupMissing:
		s = fmt.Sprintf("The source group %s [id=%d] was not found in the database.",
			e.FromGroupName,
			e.FromGroupId,
		)
		if e.ToGroupId != 0 {
			s = fmt.Sprintf("%s Link '%s' [id=%d] moved to destination group '%s' [id=%d] anyway.", s,
				e.LinkURL,
				e.LinkId,
				e.ToGroupName,
				e.ToGroupId,
			)
		}
	case shared.FromLinkGroupMissing:
		s = fmt.Sprintf("The association (aka 'link_group') between the source Group '%s' [id=%d] and the link '%s' [id=%d] was not found in the database.",
			e.FromGroupName,
			e.FromGroupId,
			e.LinkURL,
			e.LinkId,
		)
		if e.ToGroupId != 0 {
			s = fmt.Sprintf("%s Link moved to destination group '%s' [id=%d] anyway.", s,
				e.ToGroupName,
				e.ToGroupId,
			)
		}
	case shared.LinkMissing:
		s = fmt.Sprintf("The link '%s' [id=%d] was not found in the database so could not be moved.",
			e.LinkURL,
			e.LinkId,
		)
	case shared.ToGroupFound:
		s = fmt.Sprintf("The destination group '%s' [id=%d] was not found in the database so link '%s' [id=%d] could not be moved.",
			e.ToGroupName,
			e.ToGroupId,
			e.LinkURL,
			e.LinkId,
		)
	case shared.ToLinkGroupFound:
		s = fmt.Sprintf("An association (aka 'link_group') between the link '%s' [id=%d] and the destination Group '%s' [id=%d] already existed in the database so move ignored.",
			e.LinkURL,
			e.LinkId,
			e.ToGroupName,
			e.ToGroupId,
		)
	default:
		slog.Error("Unexpected link group move exception",
			"exception_type", e.Exception,
			"exception", e,
		)
	}
	return s
}
func (ee LinkGroupMoveExceptions) AsErrors() []error {
	return shared.ConvertSlice([]LinkGroupMoveException(ee), func(e LinkGroupMoveException) error {
		return e.AsError()
	})
}
func (ee LinkGroupMoveExceptions) Errors() (ss []string) {
	visited := make(map[string]struct{})
	ss = make([]string, 0, len(ee))
	for _, e := range ee {
		switch e.Exception {
		case shared.FromGroupMissing, shared.ToGroupFound:
			// Only include one for each of these
			if _, ok := visited[e.Exception]; ok {
				continue
			}
			visited[e.Exception] = struct{}{}
		case shared.FromLinkGroupMissing:
			mk := fmt.Sprintf("%d-%d", e.FromGroupId, e.LinkId)
			// Only include one for each combination
			if _, ok := visited[mk]; ok {
				continue
			}
			visited[mk] = struct{}{}
		case shared.ToLinkGroupFound:
			mk := fmt.Sprintf("%d-%d", e.ToGroupId, e.LinkId)
			// Only include one for each combination
			if _, ok := visited[mk]; ok {
				continue
			}
			visited[mk] = struct{}{}
		case shared.LinkMissing:
			if _, ok := visited[e.LinkURL]; ok {
				continue
			}
			visited[e.LinkURL] = struct{}{}
		default:
			slog.Error("Unexpected link group move exception",
				"exception_type", e.Exception,
				"exception", e,
			)
		}
		ss = append(ss, e.Error())
	}
	return ss
}
func (ee LinkGroupMoveExceptions) String() (s string) {
	return "WARNING(s):\n" + ee.Error()
}

func (ee LinkGroupMoveExceptions) Error() (s string) {
	sb := strings.Builder{}
	n := 1
	for _, e := range ee.Errors() {
		sb.WriteString(strconv.Itoa(n))
		sb.WriteString(".) ")
		sb.WriteString(e)
		sb.WriteByte('\n')
		n++
	}
	return sb.String()
}
func (ee LinkGroupMoveExceptions) AsError() error {
	return errors.Join(ee.AsErrors()...)
}

func (e LinkGroupMoveException) AsError() error {
	return errors.Join(
		ErrLinkGroupMoveException,
		fmt.Errorf("from_group=%s", e.FromGroupName),
		fmt.Errorf("from_group_id=%d", e.FromGroupId),
		fmt.Errorf("to_group=%s", e.ToGroupName),
		fmt.Errorf("to_group_id=%d", e.ToGroupId),
		fmt.Errorf("link_id=%d", e.LinkId),
		fmt.Errorf("link_url=%s", e.LinkURL),
	)
}

type MoveLinksToGroupResult struct {
	Status     LinkGroupMoveStatus
	Exceptions LinkGroupMoveExceptions
}

func (r MoveLinksToGroupResult) AsError() error {
	return r.Exceptions.AsError()
}

func (r MoveLinksToGroupResult) HasExceptions() bool {
	return len(r.Exceptions) != 0
}

func MoveLinkToGroup(ctx Context, args MoveLinkToGroupArgs) (result MoveLinksToGroupResult, err error) {
	var r storage.MoveLinksToGroupResult

	r, err = storage.MoveLinksToGroup(ctx, nil, storage.MoveLinkToGroupArgs(args))
	if err != nil {
		goto end
	}
	result = MoveLinksToGroupResult{
		Status: LinkGroupMoveStatus(r.Status),
	}
	if r.Exceptions == nil {
		goto end
	}
	result.Exceptions = shared.ConvertSlice(r.Exceptions, func(r storage.LinkGroupMoveException) LinkGroupMoveException {
		return LinkGroupMoveException(r)
	})
end:
	return result, err
}
