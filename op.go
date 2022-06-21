package goql

import "strings"

type Op int

func (o Op) Has(tgt Op) bool {
	return o&tgt == tgt
}

func (o Op) Is(tgt Op) bool {
	return o == tgt
}

//go:generate stringer -type Op -trimprefix Op
const (
	OpEq    Op = 1 << iota // =, equals, e.g. name=eq:john becomes name = 'john'
	OpNeq                  // <> or !=, not equals, e.g. name=neq:john becomes name <> 'john'
	OpLt                   // <, less than
	OpLte                  // <=, less than equals
	OpGt                   // >, greater than
	OpGte                  // >=, greater than equals
	OpLike                 // like, e.g. like:john* becomes name like 'john%' (use * instead of % in query string)
	OpIlike                // ilike, same as like, but case insensitive, e.g. name=ilike:john%
	OpIn                   // in, e.g. name=in:{1,2,3| becomes name in (1,2,3). The curly brackets is needed to indicate that the whole value is an array, since for array columns there's no way to differentiate for a single value.
	OpIs                   // is, checking for exact equality (null,true,false,unknown), e.g. age=is:null
	OpFts                  // Full-Text search using to_tsquery
	OpPlFts                // Full-Text search using plain to tsquery
	OpPhFts                // Full-Text search using phrase to tsquery
	OpWFts                 // Full-Text search using word.
	OpCs                   // @>, contains, e.g. ?tags=cs.{example,new}
	OpCd                   // <@, contained in e.g. ?values=cd.{1,2,3}
	OpOv                   // &&, Overlap
	OpSl                   // <<, strictly left of
	OpSr                   // >>, strictly right of
	OpNxr                  // &<
	OpNxl                  // &>
	OpAdj                  // -|-
	OpNot
	OpOr
	OpAnd
)

func ParseOp(op string) (Op, bool) {
	for i := OpAnd; i > 0; i = i >> 1 {
		if strings.EqualFold(i.String(), op) {
			return i, true
		}
	}

	return 0, false
}
