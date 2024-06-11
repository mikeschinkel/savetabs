package guard

import (
	"errors"
	"fmt"

	"savetabs/model"
	"savetabs/shared"
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

func ApplyDragDrop(ctx Context, args ApplyDragDropArgs) (skipped []int64, err error) {
	var parentType, dragType, dropType shared.Identifier
	var dd *shared.DragDrop

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
		skipped, err = model.MoveLinkToGroup(ctx, model.MoveLinkToGroupArgs{
			LinkIds:     args.DragIds,
			FromGroupId: args.ParentId,
			ToGroupId:   args.DropId,
		})
		if err != nil {
			me.Add(err)
		}
	}
end:
	err = me.Err()
	if err != nil {
		err = errors.Join(err, fmt.Errorf("drag_drop=%s", args.String()))
	}
	return skipped, err
}
