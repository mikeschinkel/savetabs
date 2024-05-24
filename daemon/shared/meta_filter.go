package shared

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var _ FilterItem = (*MetaFilter)(nil)

type Meta struct {
	Key   string
	Value string
}

type MetaFilter struct {
	Metas []Meta
}

func (mf MetaFilter) FilterType() *FilterType {
	return MetaFilterType
}

func (mf MetaFilter) Label() string {
	return strings.Join(mf.labels(), ",")
}

func (mf MetaFilter) Filters() []any {
	return ConvertSlice(mf.Metas, func(meta Meta) any {
		return meta
	})
}

func (g Meta) Label() string {
	return fmt.Sprintf("%s=%s", g.Key, g.Value)
}

func (mf MetaFilter) labels() []string {
	return ConvertSlice(mf.Metas, func(meta Meta) string {
		return meta.Label()
	})
}

var metaRegexp = regexp.MustCompile(`^([a-z0-_:]+)=(.+)$`)

func parseMeta(value string) (meta Meta, err error) {
	match := metaRegexp.FindStringSubmatch(value)
	if match == nil {
		err = errors.Join(
			ErrInvalidMetaFilterFormat,
			fmt.Errorf("filter_value=%s, format_expected='<meta_type>:<meta_name>'", value),
		)
		goto end
	}
	meta = Meta{
		Key:   match[1], // TODO: Perform sanitation on keys
		Value: match[2],
	}
end:
	return meta, err
}

func ParseMetaFilter(value string) (mf MetaFilter, found bool, err error) {
	var me *MultiErr
	values := strings.Split(value, ",")
	if len(values) == 0 {
		goto end
	}
	me = NewMultiErr()
	mf.Metas = make([]Meta, 0, len(values))
	for _, value := range values {
		meta, err := parseMeta(value)
		if err != nil {
			me.Add(err)
			continue
		}
		found = true
		mf.Metas = append(mf.Metas, meta)
	}
	err = me.Err()
end:
	return mf, found, err
}
