package goql

import (
	"fmt"
	"net/url"
)

type Query struct {
	Field  string
	Op     Op
	Values []string
}

func (q Query) String() string {
	return fmt.Sprintf("%s.%s:%v", q.Field, q.Op, q.Values)
}

// ParseQuery parses the query with operators, excludes
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

	return result, nil
}
