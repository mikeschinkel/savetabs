package shared

import (
	"strconv"
	"strings"
)

type Int64Slice []int64

func (ii Int64Slice) Join(separator string) string {
	return strings.Join(ii.Strings(), separator)
}

func (ii Int64Slice) Strings() []string {
	return ConvertSlice([]int64(ii), func(i int64) string {
		return strconv.FormatInt(i, 10)
	})
}
