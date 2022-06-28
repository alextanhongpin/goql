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

func (q Query) String() string {
	return fmt.Sprintf("%s.%s:%v", q.Field, q.Op, q.Values)
}

// ParseQuery parses the query with operators.
func ParseQuery(v url.Values, excludes ...string) ([]Query, error) {
	result := make([]Query, 0, len(v))

	exclude := make(map[string]bool)
	for _, val := range excludes {
		exclude[val] = true
	}

	for key, values := range v {
		// Skip zero values.
		if len(values) == 0 {
			continue
		}

		if exclude[key] {
			continue
		}

		field, operator := Split2(key, ".")
		if field == "" {
			return nil, fmt.Errorf("%w: %s", ErrUnknownField, key)
		}

		op, ok := ParseOp(operator)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownOperator, key)
		}

		result = append(result, Query{
			Field:  field,
			Op:     op,
			Values: values,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		lhs, rhs := result[i], result[j]
		byField := strings.Compare(lhs.Field, rhs.Field)
		byOp := strings.Compare(lhs.Op.String(), rhs.Op.String())

		return sortBy(+byField, +byOp)
	})

	return result, nil
}

func sortBy(dir ...int) bool {
	for _, c := range dir {
		if c != 0 {
			return c < 0
		}
	}

	return dir[len(dir)-1] < 0
}
