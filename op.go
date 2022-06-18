package goql

import "strings"

type Op int

//go:generate stringer -type Op -trimprefix Op
const (
	OpEq Op = iota
	OpNeq
	OpLt
	OpLte
	OpGt
	OpGte
	OpIn
	OpBetween
	OpIs // Literal value, e.g. now(), null
	OpIsNot

	// Range
	OpOverlap

	// MultiRange
	// Text
	// TextSearch
	// Regex

	// Time range
	OpAt
	OpContainedBy
	OpContainedAt
	OpNow
	OpToday
	OpYesterday

	// Array
	OpContains
	OpSubset
	OpSuperSet

	// JSON
)

type Ops map[Op]bool

func (o Ops) Copy() Ops {
	oo := make(Ops)

	for k, v := range o {
		oo[k] = v
	}

	return oo
}

func (o Ops) Is(v any) bool {
	switch t := v.(type) {
	case string:
		for k := range o {
			if strings.EqualFold(k.String(), t) {
				return true
			}
		}

		return false
	case int:
		return o[Op(t)]
	default:
		return false
	}
}

func (o Ops) Get(v any) (Op, bool) {
	switch t := v.(type) {
	case string:
		for k := range o {
			if strings.EqualFold(k.String(), t) {
				return k, true
			}
		}

		return 0, false
	case int:
		op := Op(t)

		return op, o[op]
	default:
		return 0, false
	}
}
