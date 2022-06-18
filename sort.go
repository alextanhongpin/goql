package goql

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	ErrInvalidSort = errors.New("goql: invalid sort")
	ErrEmptySort   = errors.New("goql: empty sort")
)

var SortKey = "sort_by"

type SortOrder string

var (
	SortOrderAscending  SortOrder = "asc"
	SortOrderDescending SortOrder = "desc"
)

type SortBy struct {
	Field string
	Order SortOrder
	Pos   int
}

func ParseSortBy(values url.Values) ([]SortBy, error) {
	if len(values[SortKey]) > 1 {
		return nil, fmt.Errorf("%w: multiple sort keys: %s", ErrInvalidSort, values[SortKey])
	}

	sortBy := strings.TrimSpace(values.Get(SortKey))
	sortables := strings.Split(sortBy, ",")

	result := make([]SortBy, len(sortables))

	for i, s := range sortables {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil, fmt.Errorf("%w: %s", ErrEmptySort, sortBy)
		}

		var o SortOrder
		switch {
		case strings.HasPrefix(s, "-"):
			s = strings.TrimPrefix(s, "-")
			o = SortOrderDescending
		case strings.HasPrefix(s, "+"):
			s = strings.TrimPrefix(s, "+")
			o = SortOrderAscending
		default:
			o = SortOrderAscending
		}

		result[i] = SortBy{
			Field: s,
			Order: o,
			Pos:   i + 1,
		}
	}

	return result, nil
}
