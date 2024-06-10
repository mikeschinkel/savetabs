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
	return fmt.Sprintf("%s:%d/%s:%s ==> %s:%s",
		args.ParentType,
		args.ParentId,
		args.DragType,
		args.DragIds,
		args.DropType,
		args.DropId,
	)
}

func ApplyDragDrop(ctx Context, args ApplyDragDropArgs) (err error) {
	var parentType, dragType, dropType shared.Identifier
	var dd *shared.DragDrop

	me := shared.NewMultiErr()
	parentType = shared.NewIdentifier(args.ParentType)
	if !parentType.Valid {
		me.Add(errors.Join(ErrInvalidIdentifier, fmt.Errorf("type=parent")))
	}
	dragType = shared.NewIdentifier(args.DragType)
	if !dragType.Valid {
		me.Add(errors.Join(ErrInvalidIdentifier, fmt.Errorf("type=drag")))
	}
	dropType = shared.NewIdentifier(args.DropType)
	if !dropType.Valid {
		me.Add(errors.Join(ErrInvalidIdentifier, fmt.Errorf("type=drop")))
	}
	dd, err = shared.DragDropByTypes(dragType, dropType)
	if err != nil {
		goto end
	}
	switch dd {
	case shared.LinkToGroupDragDrop:
		err = model.MoveLinkToGroup(ctx, model.MoveLinkToGroupArgs{
			LinkIds:     args.DragIds,
			FromGroupId: args.ParentId,
			ToGroupId:   args.DropId,
		})
	}
end:
	if err != nil {
		err = errors.Join(err, fmt.Errorf("drag_drop=%s", args.String()))
	}
	return err
}
