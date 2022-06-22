package goql

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidDirection = errors.New("goql: invalid direction")

var (
	DirectionAscending  Direction = "asc"
	DirectionDescending Direction = "desc"
	DirectionNullsFirst Direction = "nullsfirst"
	DirectionNullsLast  Direction = "nullslast"
)

type Direction string

func (o Direction) Valid() bool {
	switch o {
	case
		DirectionAscending,
		DirectionDescending,
		DirectionNullsFirst,
		DirectionNullsLast:
		return true
	default:
		return false
	}
}

type Order struct {
	Column    string
	Direction Direction
}

func NewOrder(s string) (*Order, error) {
	field, order := split2(s, ".")
	if order == "" {
		return &Order{
			Column:    field,
			Direction: DirectionAscending,
		}, nil
	}

	if o := Direction(order); !o.Valid() {
		return nil, fmt.Errorf("%w: %q", ErrInvalidDirection, order)
	}

	return &Order{
		Column:    field,
		Direction: Direction(order),
	}, nil
}

func ParseOrder(query string) ([]Order, error) {
	query = strings.TrimSpace(query)
	orders := strings.Split(query, ",")

	result := make([]Order, len(orders))

	for i, s := range orders {
		ord, err := NewOrder(s)
		if err != nil {
			return nil, err
		}

		result[i] = *ord
	}

	return result, nil
}
