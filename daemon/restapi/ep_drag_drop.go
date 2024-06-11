package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"savetabs/guard"
	"savetabs/shared"
)

type dragDropItems struct {
	Type string  `json:"type" validate:"required"`
	Ids  []int64 `json:"ids" validate:"required"`
}
type dragDropItem struct {
	Type string `json:"type" validate:"required"`
	Id   int64  `json:"id" validate:"required"`
}

type dragDrop struct {
	Parent dragDropItem  `json:"parent" validate:"required"`
	Drag   dragDropItems `json:"drag" validate:"required"`
	Drop   dragDropItem  `json:"drop" validate:"required"`
}

func (d dragDrop) String() string {
	return fmt.Sprintf("%s:%d/%s:%s => %s:%d",
		d.Parent.Type,
		d.Parent.Id,
		d.Drag.Type,
		strings.Join(shared.ConvertSlice(d.Drag.Ids, func(id int64) string {
			return strconv.FormatInt(id, 10)
		}), ","),
		d.Drop.Type,
		d.Drop.Id,
	)
}

func (a *API) PostDragDrop(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Find a better status result than "Bad Gateway"
		a.sendHTMLError(w, r, http.StatusBadGateway, err.Error())
		return
	}
	var dd dragDrop
	err = json.Unmarshal(body, &dd)
	if err != nil {
		a.sendHTMLError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	err = shared.ValidateStruct(dd)
	if err != nil {
		a.sendHTMLError(w, r, http.StatusBadRequest, err.Error())
		return
	}
	msg := dd.String()

	slog.Debug("POST name:", "drag_drop", msg)

	var skipped []int64
	skipped, err = guard.ApplyDragDrop(ctx, guard.ApplyDragDropArgs{
		ParentType: dd.Parent.Type,
		ParentId:   dd.Parent.Id,
		DragType:   dd.Drag.Type,
		DragIds:    dd.Drag.Ids,
		DropType:   dd.Drop.Type,
		DropId:     dd.Drop.Id,
	})
	if err != nil {
		a.sendHTMLError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	jr := newJSONResponse(true)
	switch {
	case len(skipped) == len(dd.Drag.Ids):
		jr.Message = fmt.Sprintf("Drop %s NOT applied; no %s IDs were associated with %s:%d",
			msg,
			dd.Drag.Type,
			dd.Parent.Type,
			dd.Parent.Id,
		)
	case len(skipped) > 0:
		jr.Message = fmt.Sprintf("Drop %s applied; Some IDs [%s] were not associated with %s:%d",
			msg,
			shared.Int64Slice(dd.Drag.Ids),
			dd.Parent.Type,
			dd.Parent.Id,
		)
	default:
		jr.Message = fmt.Sprintf("Drop Applied: %s", msg)
	}
	sendJSON(w, http.StatusOK, jr)
}
