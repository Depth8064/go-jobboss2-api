package jobboss

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	DefaultTake = 200
	MaxTake     = 200
)

// FilterOp is a JobBOSS filter operator.
type FilterOp string

const (
	FilterEqual            FilterOp = "eq"
	FilterNotEqual         FilterOp = "ne"
	FilterGreaterThan      FilterOp = "gt"
	FilterGreaterThanEqual FilterOp = "gte"
	FilterLessThan         FilterOp = "lt"
	FilterLessThanEqual    FilterOp = "lte"
	FilterIn               FilterOp = "in"
	FilterNotIn            FilterOp = "notin"
	FilterNull             FilterOp = "null"
)

// NewQuery returns an initialized url.Values map.
func NewQuery() url.Values {
	return url.Values{}
}

// ApplyPaging sets take/skip query params using JobBOSS defaults and limits.
func ApplyPaging(q url.Values, take, skip int) url.Values {
	if q == nil {
		q = url.Values{}
	}
	if take <= 0 {
		take = DefaultTake
	}
	if take > MaxTake {
		take = MaxTake
	}
	if skip < 0 {
		skip = 0
	}

	q.Set("take", fmt.Sprintf("%d", take))
	if skip > 0 {
		q.Set("skip", fmt.Sprintf("%d", skip))
	}
	return q
}

// ApplyFields sets the fields expression: fields=a,b,c.
func ApplyFields(q url.Values, fields ...string) url.Values {
	if q == nil {
		q = url.Values{}
	}
	clean := make([]string, 0, len(fields))
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if f != "" {
			clean = append(clean, f)
		}
	}
	if len(clean) > 0 {
		q.Set("fields", strings.Join(clean, ","))
	}
	return q
}

// ApplySort sets the sort expression: sort=+a,-b.
func ApplySort(q url.Values, expressions ...string) url.Values {
	if q == nil {
		q = url.Values{}
	}
	clean := make([]string, 0, len(expressions))
	for _, expr := range expressions {
		expr = strings.TrimSpace(expr)
		if expr != "" {
			clean = append(clean, expr)
		}
	}
	if len(clean) > 0 {
		q.Set("sort", strings.Join(clean, ","))
	}
	return q
}

// ApplyFilter adds a filter expression key like field[op]=value.
// If op is empty, it applies simple equality as field=value.
func ApplyFilter(q url.Values, field string, op FilterOp, value string) url.Values {
	if q == nil {
		q = url.Values{}
	}
	field = strings.TrimSpace(field)
	value = strings.TrimSpace(value)
	if field == "" || value == "" {
		return q
	}
	if strings.TrimSpace(string(op)) == "" {
		q.Set(field, value)
		return q
	}
	q.Set(fmt.Sprintf("%s[%s]", field, string(op)), value)
	return q
}

// ApplyInFilter applies an [in] filter using pipe-separated values.
func ApplyInFilter(q url.Values, field string, values ...string) url.Values {
	clean := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" {
			clean = append(clean, v)
		}
	}
	if len(clean) == 0 {
		return q
	}
	return ApplyFilter(q, field, FilterIn, strings.Join(clean, "|"))
}
