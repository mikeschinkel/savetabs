package restapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

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
	Drag   dragDropItems `json:"drag" validate:"required"`
	Drop   dragDropItem  `json:"drop" validate:"required"`
	Parent dragDropItem  `json:"parent" validate:"required"`
}

func (d dragDrop) String() string {
	return fmt.Sprintf("%s:%s => %s:%d",
		d.Drag.Type,
		strings.Join(shared.ConvertSlice(d.Drag.Ids, func(id int64) string {
			return strconv.FormatInt(id, 10)
		}), ","),
		d.Drop.Type,
		d.Drop.Id,
	)
}

func (a *API) PostDragDrop(w http.ResponseWriter, r *http.Request) {

	//ctx := context.Background()

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

	if msg == "" {
		err = ErrFailedToFixThisTODO
	}

	// TODO: Implement the following:
	//err = guard.ApplyDrop(ctx, guard.ApplyDropArgs{
	//	DragDrop: dd,
	//})

	if err != nil {
		a.sendHTMLError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	jr := newJSONResponse(true)
	jr.Message = msg
	sendJSON(w, http.StatusOK, jr)
}
