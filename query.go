package goql

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// Query represents the parsed query string with operators.
type Query struct {
	Field  string
	Op     Op
	Values []string
}

func NewQuery(query string, values []string) *Query {
	field, operator := Split2(query, ".")
	op, _ := ParseOp(operator)

	return &Query{
		Field:  field,
		Op:     op,
		Values: values,
	}
}

func (q Query) String() string {
	return fmt.Sprintf("%s.%s:%v", q.Field, q.Op, q.Values)
}

func (q *Query) Validate() error {
	if q.Field == "" {
		return fmt.Errorf("%w: %s", ErrUnknownField, q)
	}

	if !q.Op.Valid() {
		return fmt.Errorf("%w: %s", ErrUnknownOperator, q)
	}

	return nil
}

// ParseQuery parses the query with operators.
func ParseQuery(values url.Values, excludes ...string) ([]Query, error) {
	result := make([]Query, 0, len(values))

	exclude := make(map[string]bool)
	for _, val := range excludes {
		exclude[val] = true
	}

	for key, vals := range values {
		if exclude[key] || len(vals) == 0 {
			continue
		}

		query := NewQuery(key, vals)
		if err := query.Validate(); err != nil {
			return nil, err
		}

		result = append(result, *query)
	}

	SortQuery(result)

	return result, nil
}

func SortQuery(queries []Query) {
	sort.Slice(queries, func(i, j int) bool {
		lhs, rhs := queries[i], queries[j]
		byField := strings.Compare(lhs.Field, rhs.Field)
		byOp := strings.Compare(lhs.Op.String(), rhs.Op.String())

		return sortBy(+byField, +byOp)
	})
}

func sortBy(dir ...int) bool {
	for _, c := range dir {
		if c != 0 {
			return c < 0
		}
	}

	return dir[len(dir)-1] < 0
}
