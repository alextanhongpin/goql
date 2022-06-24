package goql

import "strings"

/*

By default, all datatypes are comparable: =, <>, <, <=, >, >=
Since WHERE col IN (...) is pretty common in SQL, OpIn is added for most types.
If the datatype can be NULL, then OpIs will be appended.
*/
const (
	OpsComparable = OpEq | OpNeq | OpLt | OpLte | OpGt | OpGte | OpIn | OpNotIn
	OpsNull       = OpIs | OpNot
	OpsIs         = OpIs | OpIsNot
	OpsIn         = OpIn | OpNotIn
	OpsText       = OpLike | OpNotLike | OpIlike | OpNotIlike | OpIn | OpNotIn | OpFts | OpPlFts | OpPhFts | OpWFts
	OpsRange      = OpCs | OpCd | OpOv | OpSl | OpSr | OpNxr | OpNxl | OpAdj
)

type Op int

func (o Op) Valid() bool {
	return o != 0
}

func (o Op) Has(tgt Op) bool {
	return o&tgt == tgt
}

func (o Op) Is(tgt Op) bool {
	return o == tgt
}

//go:generate stringer -type Op -trimprefix Op
const (
	OpEq       Op = 1 << iota // =, equals, e.g. name=eq:john becomes name = 'john'
	OpNeq                     // <> or !=, not equals, e.g. name=neq:john becomes name <> 'john'
	OpLt                      // <, less than
	OpLte                     // <=, less than equals
	OpGt                      // >, greater than
	OpGte                     // >=, greater than equals
	OpLike                    // like, e.g. like:john* becomes name like 'john%' (use * instead of % in query string)
	OpIlike                   // ilike, same as like, but case insensitive, e.g. name=ilike:john%
	OpNotLike                 // not like
	OpNotIlike                // not ilike
	OpIn                      // in, e.g. name=in:{1,2,3| becomes name in (1,2,3). The curly brackets is needed to indicate that the whole value is an array, since for array columns there's no way to differentiate for a single value.
	OpNotIn                   // not in
	OpIs                      // is, checking for exact equality (null,true,false,unknown), e.g. age=is:null
	OpIsNot                   // is not {null, true, false, unknown}

	// Full-Text search.
	OpFts   // Full-Text search using to_tsquery
	OpPlFts // Full-Text search using plain to tsquery
	OpPhFts // Full-Text search using phrase to tsquery
	OpWFts  // Full-Text search using word.

	// https://www.postgresql.org/docs/14/functions-range.html
	OpCs  // @>, contains, e.g. ?tags=cs:{example,new}
	OpCd  // <@, contained in e.g. ?values=cd:{1,2,3}
	OpOv  // &&, Overlap
	OpSl  // <<, strictly left of
	OpSr  // >>, strictly right of
	OpNxr // &<
	OpNxl // &>
	OpAdj // -|-

	// Conjunctions.
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

func sqlIs(unk string) bool {
	switch unk {
	case
		"0", "f", "F", "n", "no", "false", "FALSE", "False",
		"1", "t", "T", "y", "yes", "true", "TRUE", "True",
		"null", "NULL", "Null",
		"unknown", "UNKNOWN", "Unknown":
		return true
	default:
		return false
	}
}
