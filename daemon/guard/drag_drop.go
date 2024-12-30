package guard

import (
	"errors"
	"fmt"

	"savetabs/model"
	"savetabs/shared"
	"savetabs/ui"
)

type ApplyDragDropArgs struct {
	ParentType string
	ParentId   int64
	DragType   string
	DragIds    []int64
	DropType   string
	DropId     int64
}

func (args ApplyDragDropArgs) String() string {
	return fmt.Sprintf("%s:%d/%s:%s ==> %s:%d",
		args.ParentType,
		args.ParentId,
		args.DragType,
		shared.Int64Slice(args.DragIds).Join(","),
		args.DropType,
		args.DropId,
	)
}

// ApplyDragDrop applies the drag and drop request.
// TODO: Ensure errors are errors and exception messages are exceptions (rename warnings?)
func ApplyDragDrop(ctx Context, args ApplyDragDropArgs) (result MoveLinksToGroupResult, err error) {
	var parentType, dragType, dropType shared.Identifier
	var dd *shared.DragDrop
	var r model.MoveLinksToGroupResult

	me := shared.NewMultiErr()
	parentType = shared.NewIdentifier(args.ParentType)
	if !parentType.Valid {
		me.Add(errors.Join(ErrInvalidIdentifier, fmt.Errorf("type=parent, value=%s", args.ParentType)))
	}
	dragType = shared.NewIdentifier(args.DragType)
	if !dragType.Valid {
		me.Add(errors.Join(ErrInvalidIdentifier, fmt.Errorf("type=drag, value=%s", args.DragType)))
	}
	dropType = shared.NewIdentifier(args.DropType)
	if !dropType.Valid {
		me.Add(errors.Join(ErrInvalidIdentifier, fmt.Errorf("type=drop, value=%s", args.DropType)))
	}
	dd, err = shared.DragDropByTypes(dragType, dropType)
	if err != nil {
		me.Add(err)
		goto end
	}
	switch dd {
	case shared.LinkToGroupDragDrop:
		r, err = model.MoveLinkToGroup(ctx, model.MoveLinkToGroupArgs{
			LinkIds:     args.DragIds,
			FromGroupId: args.ParentId,
			ToGroupId:   args.DropId,
		})
		if err != nil {
			me.Add(err)
		}
		if r.Status.ErrorException {
			me.Add(errors.New(model.LinkGroupMoveExceptions(r.Exceptions).String()))
			goto end
		}

		result = MoveLinksToGroupResult{r}
	}
end:
	err = me.Err()
	if err != nil {
		err = errors.Join(fmt.Errorf("INFO: drag_drop=%s", args.String()), err)
	}
	return result, err
}

type MoveLinksToGroupResult struct {
	model.MoveLinksToGroupResult
}

func (r MoveLinksToGroupResult) GetExceptionsHTML(ctx Context) (HTMLResponse, error) {
	exceptions := []model.LinkGroupMoveException(r.Exceptions)
	hr, err := ui.GetExceptionsHTML(ctx, ui.ExceptionsParams{
		Exceptions: shared.ConvertSlice(exceptions, func(e model.LinkGroupMoveException) string {
			return e.Error()
		}),
	})
	return HTMLResponse{hr}, err
}
