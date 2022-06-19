package goql

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidSort = errors.New("goql: invalid sort")
	ErrEmptySort   = errors.New("goql: empty sort")
)

var SortKey = "sort_by"

var (
	SortOrderAscending  Order = "asc"
	SortOrderDescending Order = "desc"
	SortOrderNullsFirst Order = "nullsfirst"
	SortOrderNullsLast  Order = "nullslast"
)

type Order string

func (o Order) Valid() bool {
	switch o {
	case
		SortOrderAscending,
		SortOrderDescending,
		SortOrderNullsFirst,
		SortOrderNullsLast:
		return true
	default:
		return false
	}
}

type SortBy struct {
	Field string
	Order Order
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

		ext := filepath.Ext(s)
		s = s[:len(s)-len(ext)]
		ext = strings.ReplaceAll(ext, ".", "")

		o := SortOrderAscending
		if s != "" {
			if v := Order(s); v.Valid() {
				o = v
			}
		}

		result[i] = SortBy{
			Field: s,
			Order: o,
			Pos:   i + 1,
		}
	}

	return result, nil
}
