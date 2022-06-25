package goql

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidSortDirection = errors.New("goql: invalid sort direction")
	ErrInvalidSortOption    = errors.New("goql: invalid sort option")
)

//https://www.postgresql.org/docs/current/queries-order.html#:~:text=The%20NULLS%20FIRST%20and%20NULLS,order%2C%20and%20NULLS%20LAST%20otherwise.

var (
	SortDirectionAscending  SortDirection = "asc"
	SortDirectionDescending SortDirection = "desc"
	SortOptionNullsFirst    SortOption    = "nullsfirst" // Default
	SortOptionNullsLast     SortOption    = "nullslast"  // Default
)

type SortDirection string

func (o SortDirection) Valid() bool {
	return o == SortDirectionAscending || o == SortDirectionDescending
}

func (o SortDirection) DefaultOption() SortOption {
	switch o {
	case SortDirectionAscending:
		return SortOptionNullsLast
	case SortDirectionDescending:
		return SortOptionNullsFirst
	default:
		panic("goql: invalid sort direction")
	}
}

type SortOption string

func (o SortOption) Valid() bool {
	return o == SortOptionNullsFirst || o == SortOptionNullsLast
}

type Order struct {
	Field     string
	Direction SortDirection
	Option    SortOption
}

func NewOrder(s string) (*Order, error) {
	if s == "" {
		return nil, nil
	}

	field, direction, option := Split3(s, ".")
	if direction == "" {
		return &Order{
			Field:     field,
			Direction: SortDirectionAscending,
			Option:    SortDirectionAscending.DefaultOption(),
		}, nil
	}

	dir := SortDirection(direction)
	if !dir.Valid() {
		return nil, fmt.Errorf("%w: %q", ErrInvalidSortDirection, direction)
	}

	if option == "" {
		return &Order{
			Field:     field,
			Direction: dir,
			Option:    dir.DefaultOption(),
		}, nil
	}

	opt := SortOption(option)
	if !opt.Valid() {
		return nil, fmt.Errorf("%w: %q", ErrInvalidSortOption, option)
	}

	return &Order{
		Field:     field,
		Direction: dir,
		Option:    opt,
	}, nil
}

func ParseOrder(orders []string) ([]Order, error) {

	result := make([]Order, 0, len(orders))

	for _, s := range orders {
		if s == "" {
			continue
		}

		ord, err := NewOrder(s)
		if err != nil {
			return nil, err
		}

		result = append(result, *ord)
	}

	return result, nil
}
