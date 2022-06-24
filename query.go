package goql

import (
	"fmt"
	"net/url"
)

type Query struct {
	Field string
	Op    Op
	Value string
}

func (q Query) String() string {
	return fmt.Sprintf("%s=%s:%s", q.Field, q.Op, q.Value)
}

// ParseQuery parses the query with operators, excluding
// AND and OR conjunctions.
func ParseQuery(v url.Values) []Query {
	result := make([]Query, 0, len(v))

	for field, params := range v {
		switch field {
		case "and", "or":
			continue
		}

		for _, param := range params {
			key, val := split2(param, ":")
			if val == "" {
				continue
			}

			op, ok := ParseOp(key)
			if !ok {
				continue
			}

			result = append(result, Query{
				Field: field,
				Op:    op,
				Value: val,
			})
		}
	}

	return result
}
