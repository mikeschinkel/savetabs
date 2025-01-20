package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/safehtml"
	"github.com/mikeschinkel/savetabs/daemon/guard"
	"github.com/mikeschinkel/savetabs/daemon/shared"
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

func (a *API) PostHtmlDragDrop(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	result, dd, err := a.dragDrop(ctx, w, r)
	switch {
	case err != nil:
		a.sendHTMLError(w, r, http.StatusBadRequest, err.Error())
	case result.HasExceptions():
		var hr guard.HTMLResponse
		hr, err = result.GetExceptionsHTML(ctx)
		if err != nil {
			hr.StatusCode = http.StatusInternalServerError
			//goland:noinspection GoDfaErrorMayBeNotNil
			hr.HTML = safehtml.HTMLConcat(hr.HTML, shared.MakeSafeHTML(err.Error()))
		}
		//goland:noinspection GoDfaErrorMayBeNotNil
		a.sendHTMLStatus(w, hr.StatusCode, hr.HTML)
	default:
		a.sendHTMLStatus(w, http.StatusNoContent, shared.MakeSafeHTML(dd.String()))
	}
}

func (a *API) PostDragDrop(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	result, dd, err := a.dragDrop(ctx, w, r)
	jr := newJSONResponse(true)
	switch {
	case err != nil:
		jr.Message = err.Error()
		jr.HTTPStatus = http.StatusBadRequest //TODO: Fix to acknowedge 5xx errors
	case result.HasExceptions():
		jr.Message = result.Exceptions.String()
		jr.HTTPStatus = http.StatusAccepted
	default:
		jr.Message = fmt.Sprintf("Drop Applied: %s", dd.String())
	}
	sendJSON(w, jr.HTTPStatus, jr)
}

func (a *API) dragDrop(ctx Context, w http.ResponseWriter, r *http.Request) (result guard.MoveLinksToGroupResult, dd dragDrop, err error) {
	var msg string

	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.Join(ErrReadingHTTPBody, err)
		goto end
	}
	err = json.Unmarshal(body, &dd)
	if err != nil {
		err = errors.Join(ErrUnmarshallingJSON, err)
		goto end
	}
	err = shared.ValidateStruct(dd)
	if err != nil {
		err = errors.Join(ErrValidatingHTTPRequest, err)
		goto end
	}
	msg = dd.String()

	slog.Debug("POST name:", "drag_drop", msg)

	result, err = guard.ApplyDragDrop(ctx, guard.ApplyDragDropArgs{
		ParentType: dd.Parent.Type,
		ParentId:   dd.Parent.Id,
		DragType:   dd.Drag.Type,
		DragIds:    dd.Drag.Ids,
		DropType:   dd.Drop.Type,
		DropId:     dd.Drop.Id,
	})
end:
	return result, dd, err
}
