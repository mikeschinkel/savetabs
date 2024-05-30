package guard

import (
	"savetabs/model"
)

func LoadGroupName(ctx Context, groupId int64) (name string, err error) {
	return model.LoadGroupName(ctx, groupId)
}
